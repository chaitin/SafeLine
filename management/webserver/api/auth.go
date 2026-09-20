package api

import (
	"crypto/subtle"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

// BootstrapWindowEnv names the environment variable that holds how long the
// anonymous TFA bootstrap stays open after it was opened. It accepts a Go
// duration, and "0" or "off" leaves it open until the first login.
//
// The window is what keeps an installation that nobody ever bound from being
// claimable forever: the account can be bound without a credential while the
// window is open, and `mgt -reset_user admin` opens it again, which is what an
// operator does when they want to bind a console that was left alone.
const BootstrapWindowEnv = "MGT_BOOTSTRAP_WINDOW"

// DefaultBootstrapWindow is how long the anonymous TFA bootstrap stays open
// when BootstrapWindowEnv is not set. It is long enough for the operator of an
// installation that is being set up to scan the code, and short enough that an
// installation which is left unbound does not stay open to the network.
const DefaultBootstrapWindow = 30 * time.Minute

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

	now := time.Now()
	usedStep, valid := verifyPasscode(user.TFASecret, params.Passcode, user.TFALastUsedStep, now)
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

	// The passcode is spent here, in one conditional update, and the session
	// below is only created when this request is the one that spent it. A plain
	// Save after the check would leave a window in which two requests that
	// carry the same passcode both read the same stored time step and both
	// accept it.
	consumed, err := consumePasscode(db, &user, usedStep, now.Unix()/int64(OtpOpts.Period), now)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when creating sessions"}, http.StatusInternalServerError)
		return
	}
	if !consumed {
		loginThrottle.fail(keys)
		logger.Warnf("Refused a login from %s: the passcode of time step %d was already used", c.ClientIP(), usedStep)
		response.Error(c, response.JSONBody{Err: response.ErrWrongPasscode, Msg: "Failed to verify your passcode"}, http.StatusUnauthorized)
		return
	}

	loginThrottle.reset(keys)

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

// checkBootstrapWindow returns an error when the window in which the account
// may be bound without a credential has closed.
//
// A passcode is the only credential of the account, so the secret that the
// caller of GET /api/OTPUrl receives is the account. Leaving the endpoint open
// until somebody happens to log in means an installation that is deployed and
// then forgotten can be claimed by whoever reaches the console port first; the
// window bounds that to the time around the setup, and every console that is
// bound afterwards reaches this endpoint only to be told to run
// `mgt -reset_user`.
func checkBootstrapWindow() error {
	window := BootstrapWindow()
	if window <= 0 {
		return nil
	}

	openedAt, ok := model.BootstrapOpenedAt()
	if !ok {
		// The record is written when the installation is set up, so a missing
		// one means this database was not written by this release. The account
		// is not left bound to nobody for that reason: the window starts being
		// measured from the moment the account was created.
		var user model.User
		database.GetDB().Where(&model.User{Username: constants.SuperUser}).First(&user)
		if openedAt, ok = user.CreatedAt, !user.CreatedAt.IsZero(); !ok {
			openedAt = time.Now()
		}

		logger.Warnf("No record of when the TFA bootstrap was opened, measuring the window of %s from %s", window, openedAt.Format(time.RFC3339))
	}

	if elapsed := time.Since(openedAt); elapsed > window {
		return fmt.Errorf("the anonymous TFA bootstrap closed %s after it was opened at %s", window, openedAt.Format(time.RFC3339))
	}

	return nil
}

// BootstrapWindow returns how long the anonymous TFA bootstrap stays open, and
// zero when it stays open until the first login.
func BootstrapWindow() time.Duration {
	value := strings.TrimSpace(os.Getenv(BootstrapWindowEnv))
	if value == "" {
		return DefaultBootstrapWindow
	}

	if strings.EqualFold(value, "off") {
		return 0
	}

	window, err := time.ParseDuration(value)
	if err != nil || window < 0 {
		logger.Warnf("Ignoring %s=%q: it is not a duration, keeping the window of %s", BootstrapWindowEnv, value, DefaultBootstrapWindow)
		return DefaultBootstrapWindow
	}

	return window
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
			//
			// The window applies here as well, because adopting the pending
			// secret is the same thing as being handed one.
			if err := checkBootstrapWindow(); err != nil {
				logger.Warnf("Refused to hand out the pending TFA secret to %s: %s", clientIP, err)
				response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "TFA bootstrap is not allowed for this request"}, http.StatusForbidden)
				return
			}

			adopted, err := adoptBootstrapSecret(db, user.TFASecret, clientIP)
			if err != nil {
				logger.Error(err)
				response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating otp qrcode"}, http.StatusInternalServerError)
				return
			}
			if !adopted {
				refusePendingSecret(c, clientIP)
				return
			}
		} else if user.TFABootstrapIP != clientIP {
			refusePendingSecret(c, clientIP)
			return
		}

		logger.Debugf("Reusing the pending TFA secret for %s", clientIP)
		response.Success(c, gin.H{"url": otpURL(user.TFASecret)})
		return
	}

	if err := checkBootstrapWindow(); err != nil {
		logger.Warnf("Refused to issue a TFA secret to %s: %s; run 'mgt -reset_user admin' to open the bootstrap again", clientIP, err)
		response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "The TFA bootstrap of this installation is closed, run 'mgt -reset_user admin' to bind a new authenticator"}, http.StatusForbidden)
		return
	}

	otpKey, err := totp.Generate(OtpOpts)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating otp qrcode"}, http.StatusInternalServerError)
		return
	}

	// The secret is installed by an update that only matches while the row
	// still has none, so two calls that arrive together cannot both be told
	// that the secret they generated is the one the account carries: the
	// caller that loses reads the stored secret below and hands that one out
	// when it belongs to the same address, which is what a console that
	// retries produces, and is refused otherwise.
	issued, err := claimBootstrapSecret(db, clientIP, otpKey.Secret())
	if err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating otp qrcode"}, http.StatusInternalServerError)
		return
	}
	if !issued {
		logger.Warnf("Lost the TFA bootstrap race to another caller, %s is served the stored secret if it is its own", clientIP)
		var stored model.User
		db.Where(&model.User{Username: constants.SuperUser}).First(&stored)
		if stored.TFASecret == "" || stored.TFABootstrapIP != clientIP {
			refusePendingSecret(c, clientIP)
			return
		}

		response.Success(c, gin.H{"url": otpURL(stored.TFASecret)})
		return
	}

	logger.Warnf("Issued a new TFA secret to %s; it can only be bound and used from that address until the first login", clientIP)
	response.Success(c, gin.H{"url": otpKey.URL()})
}

// refusePendingSecret answers a caller that may not be handed the pending TFA
// secret, because it was issued to another address.
func refusePendingSecret(c *gin.Context, clientIP string) {
	logger.Warnf("Refused to hand out the pending TFA secret to %s: it belongs to another address", clientIP)
	response.Error(c, response.JSONBody{Err: response.ErrLoginRequired, Msg: "The pending TFA secret belongs to another address, run 'mgt -reset_user <user>' to bind a new one"}, http.StatusForbidden)
}

// claimBootstrapSecret installs secret as the pending TFA secret of the
// account and reports whether this call is the one that installed it.
//
// The update matches only while the row carries no secret yet, which is what
// makes the claim exclusive. Reading the row and writing it back would let two
// callers that arrive while the account has no secret both store the one they
// generated, so both would be told the account is bound to the secret they
// hold, while the database keeps only the last write.
func claimBootstrapSecret(db *database.PostgresDB, clientIP, secret string) (bool, error) {
	result := db.Model(&model.User{}).
		Where("username = ? AND (tfa_secret = '' OR tfa_secret IS NULL)", constants.SuperUser).
		Updates(map[string]interface{}{
			"tfa_secret":       secret,
			"tfa_bootstrap_ip": clientIP,
		})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}

// adoptBootstrapSecret records clientIP as the address of a pending secret
// that was stored before the address was recorded.
//
// The update is restricted to the secret that was read, so a caller that lost a
// race against the installation of a new secret adopts nothing.
func adoptBootstrapSecret(db *database.PostgresDB, pending, clientIP string) (bool, error) {
	result := db.Model(&model.User{}).
		Where("username = ? AND tfa_secret = ? AND (tfa_bootstrap_ip = '' OR tfa_bootstrap_ip IS NULL)", constants.SuperUser, pending).
		Updates(map[string]interface{}{
			"tfa_bootstrap_ip": clientIP,
		})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}

// consumePasscode records an accepted passcode together with the login it
// belongs to, as a single conditional update.
//
// The update is what decides whether the time step of the passcode was already
// spent: it only matches while the stored step is below the accepted one, so
// the second of two concurrent requests that carry the same passcode affects no
// row and is treated as a replay. currentStep is the step the passcode was
// verified against; a stored step ahead of it means the clock moved backwards,
// which this process treats as stale (see verifyPasscode) and repairs here.
func consumePasscode(db *database.PostgresDB, user *model.User, usedStep, currentStep int64, now time.Time) (bool, error) {
	result := db.Model(&model.User{}).
		Where("id = ? AND (tfa_last_used_step < ? OR tfa_last_used_step > ?)", user.ID, usedStep, currentStep+1).
		Updates(map[string]interface{}{
			"tfa_last_used_step": usedStep,
			// The bootstrap is over: the secret is bound, so it is no longer
			// pending and no longer tied to the address that fetched it.
			"tfa_bootstrap_ip": "",
			"last_login_time":  now.Unix(),
			"is_enabled":       true,
		})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
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
