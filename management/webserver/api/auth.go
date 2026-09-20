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

// BootstrapTokenHeader is the header that carries the value of
// BootstrapTokenEnv.
const BootstrapTokenHeader = "X-Bootstrap-Token"

const (
	// bootstrapFailureThreshold is how many rejected bootstrap attempts a
	// client address may make before it has to wait.
	bootstrapFailureThreshold = 3

	// bootstrapFailureWindow is how long a rejected attempt is remembered.
	bootstrapFailureWindow = 10 * time.Minute

	// bootstrapMaxLockout caps the wait imposed on a caller that keeps trying.
	bootstrapMaxLockout = 5 * time.Minute
)

// bootstrapThrottle slows down guessing of the bootstrap token.
//
// GET /api/OTPUrl is reachable without a session, and the token comparison is
// the only thing between a caller and the TFA secret of the account, so a
// caller that may try values without limit could sit on the endpoint until it
// guesses one. The keys are per client address: a key shared by every caller
// would hand an attacker a way to lock the administrator out of the binding
// instead.
var bootstrapThrottle = newThrottle(bootstrapFailureThreshold, bootstrapFailureWindow, bootstrapMaxLockout)

// bootstrapKeys returns the throttle keys of a bootstrap attempt.
func bootstrapKeys(clientIP string) []string {
	if clientIP == "" {
		return nil
	}

	return []string{fmt.Sprintf("bootstrap:%s", clientIP)}
}

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

	// Until the first successful login the account is still bound to the address
	// the TFA secret was handed to. Without this check a caller that reached
	// GET /api/OTPUrl first would hold a secret the administrator is about to
	// bind, and knowing the secret is enough to log in.
	if user.LastLoginTime == 0 && user.TFASecret != "" && user.TFABootstrapIP != "" &&
		user.TFABootstrapIP != c.ClientIP() {
		loginThrottle.fail(keys)
		logger.Warnf("Refused a login from %s before the first successful one: the pending TFA secret was issued to %s", c.ClientIP(), user.TFABootstrapIP)
		response.Error(c, response.JSONBody{Err: response.ErrWrongPasscode, Msg: "Failed to verify your passcode"}, http.StatusUnauthorized)
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
	// The bootstrap is over: the secret is bound, so it is no longer pending
	// and no longer tied to the address that fetched it.
	user.TFABootstrapIP = ""
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

	clientIP := c.ClientIP()
	throttleKeys := bootstrapKeys(clientIP)
	if key, wait := bootstrapThrottle.lockout(throttleKeys); wait > 0 {
		logger.Warnf("TFA bootstrap throttled by %s, retry in %s", key, wait)
		c.Header("Retry-After", strconv.FormatInt(int64(wait.Seconds())+1, 10))
		response.Error(c, response.JSONBody{Err: response.ErrTooManyAttempts, Msg: "Too many attempts, please try again later"}, http.StatusTooManyRequests)
		return
	}

	if err := checkBootstrapToken(c); err != nil {
		bootstrapThrottle.fail(throttleKeys)
		logger.Warnf("Refused to hand out the TFA secret to %s: %s", clientIP, err)
		response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "TFA bootstrap is not allowed for this request"}, http.StatusUnauthorized)
		return
	}

	bootstrapThrottle.reset(throttleKeys)

	// The console displays the QR code on the login page, before any session
	// exists, so this endpoint is reachable while the account has never logged
	// in. It must not hand out a fresh secret on every call, and the secret it
	// does hand out belongs to the one address that asked for it: an anonymous
	// caller that fetched a secret the administrator is about to bind would
	// otherwise be able to log in with it, because the passcode is the only
	// credential.
	if user.TFASecret != "" {
		if user.TFABootstrapIP == "" {
			// An installation that already carries a pending secret from a
			// release that did not record the address keeps working: the first
			// caller adopts it, exactly as it did before.
			user.TFABootstrapIP = clientIP
			db.Save(&user)
		} else if user.TFABootstrapIP != clientIP {
			logger.Warnf("Refused to hand out the pending TFA secret to %s: it was issued to %s", clientIP, user.TFABootstrapIP)
			response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "The pending TFA secret belongs to another address, run 'mgt -reset_user <user>' to bind a new one"}, http.StatusForbidden)
			return
		}

		logger.Debugf("Reusing the pending TFA secret for %s", clientIP)
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
	user.TFABootstrapIP = clientIP
	db.Save(&user)
	logger.Warnf("Issued a new TFA secret to %s; it can only be bound and used from that address until the first login", clientIP)
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
//
// The empty value keeps the bootstrap open on purpose: the console that ships
// in the management image asks this endpoint for the secret without knowing it,
// so requiring a token here would leave a fresh installation without any way to
// bind the account. That trade off is announced at startup (see main.go) so
// that it is a visible choice rather than a silent default.
func checkBootstrapToken(c *gin.Context) error {
	expected := os.Getenv(BootstrapTokenEnv)
	if expected == "" {
		return nil
	}

	if subtle.ConstantTimeCompare([]byte(c.GetHeader(BootstrapTokenHeader)), []byte(expected)) == 1 {
		return nil
	}

	return fmt.Errorf("missing or wrong %s header", BootstrapTokenHeader)
}

func GetUser(c *gin.Context) {
	db := database.GetDB()
	user := model.User{Username: constants.SuperUser}
	db.First(&user)
	response.Success(c, gin.H{"id": user.ID, "username": user.Username})
}
