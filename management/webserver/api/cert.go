package api

import (
	"crypto/x509/pkix"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

type postSSLCertRequest struct {
	Hostname string `json:"hostname"`
}

const (
	// certGenerationThreshold is how many key pairs a client address may ask
	// for within certGenerationWindow before it has to wait.
	certGenerationThreshold = 5

	// certGenerationWindow is how long a generation is counted.
	certGenerationWindow = 10 * time.Minute

	// certGenerationMaxLockout caps the wait imposed on a caller that keeps
	// asking.
	certGenerationMaxLockout = 5 * time.Minute

	// maxUploadedCertBytes is the largest certificate or key this endpoint will
	// write. A PEM certificate or key is far smaller than this.
	maxUploadedCertBytes int64 = 1 << 20
)

// certGenerationThrottle bounds how often this endpoint may spend a key
// generation.
//
// Every accepted request generates a fresh 4096 bit RSA key in the request
// goroutine; the random file name prefix makes the existence check in
// WriteCertIfNotExist a no-op, so nothing is ever reused. An authenticated
// session that loops over this route therefore keeps the process busy with RSA
// key generation and fills the certificate directory of the site, which is the
// same process that compiles and pushes the detection rules. The throttle is
// counted per client address, and one generation is counted whether it
// succeeds or not.
var certGenerationThrottle = newThrottle(certGenerationThreshold, certGenerationWindow, certGenerationMaxLockout)

// certUploadThrottle bounds how often an authenticated session may write into
// the certificate directory. Generation already has a throttle; upload did not.
var certUploadThrottle = newThrottle(certGenerationThreshold, certGenerationWindow, certGenerationMaxLockout)

// certGenerationKeys returns the throttle keys of a certificate request.
func certGenerationKeys(clientIP string) []string {
	if clientIP == "" {
		return nil
	}

	return []string{fmt.Sprintf("sslcert:%s", clientIP)}
}

// SSLCertDir is the dir of tengine conf, not mgt-api nginx certs dir defined by constants.CertsPath
const (
	CRT = ".crt"
	PEM = ".pem"
	KEY = ".key"

	SSLCertDir = "certs"
)

func PostUploadSSLCert(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	// The console uploads both the certificate and the private key through
	// this endpoint. The extension decides which of the two it is, and the key
	// must not be created with the default (0666 before umask) mode.
	mode := os.FileMode(0644)
	switch filepath.Ext(file.Filename) {
	case CRT:
		logger.Debugf("File: %v is valid", file.Filename)
	case PEM, KEY:
		logger.Debugf("File: %v is valid", file.Filename)
		mode = os.FileMode(0600)
	default:
		logger.Errorf("Filename: %s, ext: %s", file.Filename, filepath.Ext(file.Filename))
		response.Error(c, response.JSONBody{Err: response.ErrWrongFileType, Msg: "Wrong file type, please upload a file of .crt or .key"}, http.StatusUnsupportedMediaType)
		return
	}

	if file.Size <= 0 || file.Size > maxUploadedCertBytes {
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Certificate file is empty or too large"}, http.StatusBadRequest)
		return
	}

	clientIP := c.ClientIP()
	throttleKeys := certGenerationKeys(clientIP)
	if key, wait := certUploadThrottle.lockout(throttleKeys); wait > 0 {
		logger.Warnf("Certificate upload throttled by %s, retry in %s", key, wait)
		c.Header("Retry-After", strconv.FormatInt(int64(wait.Seconds())+1, 10))
		response.Error(c, response.JSONBody{Err: response.ErrTooManyAttempts, Msg: "Too many certificate uploads, please try again later"}, http.StatusTooManyRequests)
		return
	}
	certUploadThrottle.fail(throttleKeys)

	var dstPath string
	filename := fmt.Sprintf("%s_%s", utils.RandStr(16), filepath.Base(file.Filename))
	if config.GlobalConfig.Server.DevMode {
		dstPath = filepath.Join("./nginx", SSLCertDir, filename)
	} else {
		dstPath = filepath.Join(config.GlobalConfig.NgxResDir, SSLCertDir, filename)
	}
	if err = saveUploadedFile(file, dstPath, mode); err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when saving file"}, http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"filename": filename})
}

// saveUploadedFile stores an uploaded file with an explicit permission mode.
//
// gin.Context.SaveUploadedFile calls os.Create, which asks for 0666 (0644
// after the usual umask) for every file it is given, including private keys.
// O_CREATE|O_EXCL also makes sure an existing file - or a symlink pointing
// outside the certificate directory - is never followed.
func saveUploadedFile(file *multipart.FileHeader, dstPath string, mode os.FileMode) error {
	if err := utils.EnsureFileDir(dstPath); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, io.LimitReader(src, maxUploadedCertBytes)); err != nil {
		_ = dst.Close()
		return err
	}

	return dst.Close()
}

func PostSSLCert(c *gin.Context) {
	var params postSSLCertRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	// The host name is written into the certificate as a DNS name, and a value
	// that is not a host name produces a certificate that stands for nothing in
	// particular. It is checked with the same rule the front end uses for a
	// server name before it reaches the certificate.
	if params.Hostname == "" {
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "The host name of the certificate is required"}, http.StatusBadRequest)
		return
	}

	if err := model.ValidateServerName(params.Hostname); err != nil {
		logger.Warn(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: err.Error()}, http.StatusBadRequest)
		return
	}

	clientIP := c.ClientIP()
	throttleKeys := certGenerationKeys(clientIP)
	if key, wait := certGenerationThrottle.lockout(throttleKeys); wait > 0 {
		logger.Warnf("Certificate generation throttled by %s, retry in %s", key, wait)
		c.Header("Retry-After", strconv.FormatInt(int64(wait.Seconds())+1, 10))
		response.Error(c, response.JSONBody{Err: response.ErrTooManyAttempts, Msg: "Too many certificate requests, please try again later"}, http.StatusTooManyRequests)
		return
	}

	// The generation below is the expensive part of this handler, so the
	// attempt is recorded up front: a caller that only ever produces failing
	// requests must not get an unlimited number of key generations either.
	certGenerationThrottle.fail(throttleKeys)

	filePrefix := utils.RandStr(16)
	certFilename := fmt.Sprintf("%s_backend.crt", filePrefix)
	keyFilename := fmt.Sprintf("%s_backend.key", filePrefix)

	var certPath, keyPath string
	if config.GlobalConfig.Server.DevMode {
		certPath = filepath.Join("./management", SSLCertDir, certFilename)
		keyPath = filepath.Join("./management", SSLCertDir, keyFilename)
	} else {
		certPath = filepath.Join(config.GlobalConfig.NgxResDir, SSLCertDir, certFilename)
		keyPath = filepath.Join(config.GlobalConfig.NgxResDir, SSLCertDir, keyFilename)
	}
	if err := utils.WriteCertIfNotExist(
		certPath,
		keyPath,
		func() ([]byte, []byte, error) {
			return utils.GenerateCert(
				[]string{params.Hostname},
				3650,
				4096,
				&pkix.Name{
					Country:            []string{},
					Province:           []string{},
					Locality:           []string{},
					Organization:       []string{},
					OrganizationalUnit: []string{},
					CommonName:         "",
				},
				false,
			)
		}); err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when generating certs"}, http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"crt": certFilename, "key": keyFilename})
}
