package api

import (
	"crypto/x509/pkix"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

type postSSLCertRequest struct {
	Hostname string `json:"hostname"`
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
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()
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
