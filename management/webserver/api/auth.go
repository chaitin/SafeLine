package api

import (
	"math"
	"net/http"
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
)

var logger = log.GetLogger("api")

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

	valid := totp.Validate(params.Passcode, user.TFASecret)
	if !valid {
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

	user.LastLoginTime = time.Now().Unix()
	user.IsEnabled = true
	db.Save(&user)

	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: 3600 * 24 * 7,
		//Domain: options.Domain,
		//HttpOnly: true,
		//SameSite: http.SameSiteLaxMode,
		//Secure:   false,
	})
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
	otpKey, err := totp.Generate(OtpOpts)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating otp qrcode"}, http.StatusInternalServerError)
		return
	}

	db := database.GetDB()

	// only SuperUser in v0.9
	user := model.User{Username: constants.SuperUser}
	db.First(&user)
	if user.LastLoginTime > 0 {
		// already bind tfa, because tfa binding is mandatory when login.
		response.Success(c, gin.H{"url": ""})
		return
	}

	user.TFASecret = otpKey.Secret()
	db.Save(&user)
	response.Success(c, gin.H{"url": otpKey.URL()})
}

func GetUser(c *gin.Context) {
	db := database.GetDB()
	user := model.User{Username: constants.SuperUser}
	db.First(&user)
	response.Success(c, gin.H{"id": user.ID, "username": user.Username})
}
