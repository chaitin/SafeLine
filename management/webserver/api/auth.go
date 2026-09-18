package api

import (
	"crypto/subtle"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/sessionopts"
)

var logger = log.GetLogger("api")

// loginThrottle slows down passcode guessing. The management console has a
// single account whose only credential is a six digit passcode.
var loginThrottle = newLoginGuard()

// BootstrapTokenEnv names the environment variable that gates the anonymous
// TFA bootstrap. See GetOTPUrl.
const BootstrapTokenEnv = "MGT_BOOTSTRAP_TOKEN"

// bootstrapTokenHeader is the header that carries the value of
// BootstrapTokenEnv.
const bootstrapTokenHeader = "X-Bootstrap-Token"

var OtpOpts = totp.GenerateOpts{
	Issuer:      constants.ProductName,
	AccountName: constants.SuperUser,
	Period:      30, // seconds
	Digits:      otp.DigitsSix,
	Algorithm:   otp.AlgorithmSHA1,
}

type PostLoginRequest struct {
	Passcode  string `json:"passcode"`
	Timestamp int64  `json:"timestamp"`
}

func PostLogin(c *gin.Context) {
	var params PostLoginRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()

	// only SuperUser in v0.9
	var user model.User
	db.Where(&model.User{Username: constants.SuperUser}).First(&user)

	keys := loginKeys(c.ClientIP(), user.Username)
	if key, wait := loginThrottle.lockout(keys); wait > 0 {
		logger.Warnf("Login throttled by %s, retry in %s", key, wait)
		c.Header("Retry-After", strconv.FormatInt(int64(wait.Seconds())+1, 10))
		response.Error(c, response.JSONBody{Err: response.ErrTooManyAttempts, Msg: "Too many failed login attempts, please try again later"}, http.StatusTooManyRequests)
		return
	}

	usedStep, valid := verifyPasscode(user.TFASecret, params.Passcode, user.TFALastUsedStep, time.Now())
	if !valid {
		loginThrottle.fail(keys)
		logger.Warnf("Failed login attempt from %s", c.ClientIP())

		millisecondTimestamp := params.Timestamp
		localTimeStamp := time.Now()
		if millisecondTimestamp > 0 {
			logger.Debugf("will valid otp frontend timestamp:%v, local timestamp:%v", millisecondTimestamp, localTimeStamp)
			otpTime := time.Unix(millisecondTimestamp/1000, (millisecondTimestamp%1000)*int64(time.Millisecond))
			timeSub := localTimeStamp.Sub(otpTime)
			seconds := math.Abs(timeSub.Seconds())
			if seconds >= 60 {
				logger.Errorf("otp timestamp gap is more than a minute")
				response.Error(c, response.JSONBody{Err: response.ErrWrongTimeGap, Msg: "otp timestamp gap is more than a minute"}, http.StatusUnauthorized)
				return
			}
		}
		response.Error(c, response.JSONBody{Err: response.ErrWrongPasscode, Msg: "Failed to verify your passcode"}, http.StatusUnauthorized)
		return
	}

	loginThrottle.reset(keys)

	user.LastLoginTime = time.Now().Unix()
	user.IsEnabled = true
	// Remember the time step of the accepted passcode so that the same one
	// cannot be used a second time.
	user.TFALastUsedStep = usedStep
	db.Save(&user)

	session := sessions.Default(c)
	session.Options(sessionopts.Options(c))
	session.Set(constants.DefaultSessionUserKey, user.ID)
	if err := session.Save(); err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when creating sessions"}, http.StatusInternalServerError)
		return
	}

	response.Success(c, nil)
}

func PostLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()

	if err := session.Save(); err != nil {
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when creating sessions"}, http.StatusInternalServerError)
		return
	}
	response.Success(c, nil)
}

func GetOTPUrl(c *gin.Context) {
	db := database.GetDB()

	// only SuperUser in v0.9
	user := model.User{Username: constants.SuperUser}
	db.First(&user)
	if user.LastLoginTime > 0 {
		// already bind tfa, because tfa binding is mandatory when login.
		response.Success(c, gin.H{"url": ""})
		return
	}

	if err := checkBootstrapToken(c); err != nil {
		logger.Warnf("Refused to hand out the TFA secret to %s: %s", c.ClientIP(), err)
		response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "TFA bootstrap is not allowed for this request"}, http.StatusUnauthorized)
		return
	}

	// The console displays the QR code on the login page, before any session
	// exists, so this endpoint is reachable while the account has never logged
	// in. It must not hand out a fresh secret on every call: an anonymous
	// caller that replaced a secret the administrator had already scanned would
	// take over the account, because the passcode is the only credential.
	if user.TFASecret != "" {
		logger.Debugf("Reusing the pending TFA secret for %s", c.ClientIP())
		response.Success(c, gin.H{"url": otpURL(user.TFASecret)})
		return
	}

	otpKey, err := totp.Generate(OtpOpts)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating otp qrcode"}, http.StatusInternalServerError)
		return
	}

	user.TFASecret = otpKey.Secret()
	db.Save(&user)
	logger.Warnf("Issued a new TFA secret to %s; the administrator has to bind it before the first login", c.ClientIP())
	response.Success(c, gin.H{"url": otpKey.URL()})
}

// verifyPasscode checks a TOTP passcode and returns the time step it belongs to.
//
// lastUsedStep is the step of the passcode that was accepted last; every step
// up to it is refused, which is what makes a captured passcode unusable: the
// plain totp.Validate call the handler used before accepted the same code for
// as long as it stayed inside the window of the current step.
func verifyPasscode(secret, passcode string, lastUsedStep int64, now time.Time) (int64, bool) {
	if secret == "" || passcode == "" {
		return 0, false
	}

	period := int64(OtpOpts.Period)
	currentStep := now.Unix() / period
	if lastUsedStep > currentStep+1 {
		// The recorded step is ahead of the current window, which means the
		// clock moved backwards. Keeping it would lock the account out until
		// the clock catches up, so the record is treated as stale.
		lastUsedStep = 0
	}

	opts := totp.ValidateOpts{
		Period:    OtpOpts.Period,
		Digits:    OtpOpts.Digits,
		Algorithm: OtpOpts.Algorithm,
	}

	// The current step and its two neighbours: the same window totp.Validate
	// accepts by default, so a client clock that is off by one period still
	// works.
	var matched int64
	for _, step := range []int64{currentStep - 1, currentStep, currentStep + 1} {
		if step <= lastUsedStep {
			continue
		}

		code, err := totp.GenerateCodeCustom(secret, time.Unix(step*period, 0), opts)
		if err != nil {
			continue
		}

		if subtle.ConstantTimeCompare([]byte(code), []byte(passcode)) == 1 {
			matched = step
		}
	}

	return matched, matched != 0
}

// otpURL rebuilds the otpauth URL that totp.Generate returns for a secret.
func otpURL(secret string) string {
	values := url.Values{}
	values.Set("secret", secret)
	values.Set("issuer", OtpOpts.Issuer)
	values.Set("period", strconv.FormatUint(uint64(OtpOpts.Period), 10))
	values.Set("algorithm", OtpOpts.Algorithm.String())
	values.Set("digits", OtpOpts.Digits.String())

	bootstrapURL := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + OtpOpts.Issuer + ":" + OtpOpts.AccountName,
		RawQuery: values.Encode(),
	}

	return bootstrapURL.String()
}

// checkBootstrapToken enforces the optional out of band token of the anonymous
// TFA bootstrap.
//
// Installations that do not want the account to be bindable by whoever reaches
// the API first set MGT_BOOTSTRAP_TOKEN and deliver it to the administrator
// through a channel of their choice (the installation output, for example).
func checkBootstrapToken(c *gin.Context) error {
	expected := os.Getenv(BootstrapTokenEnv)
	if expected == "" {
		return nil
	}

	if subtle.ConstantTimeCompare([]byte(c.GetHeader(bootstrapTokenHeader)), []byte(expected)) == 1 {
		return nil
	}

	return fmt.Errorf("missing or wrong %s header", bootstrapTokenHeader)
}

func GetUser(c *gin.Context) {
	db := database.GetDB()
	user := model.User{Username: constants.SuperUser}
	db.First(&user)
	response.Success(c, gin.H{"id": user.ID, "username": user.Username})
}
