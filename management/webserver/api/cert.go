package api

import (
	"crypto/x509/pkix"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
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
	switch filepath.Ext(file.Filename) {
	case CRT, PEM, KEY:
		logger.Debugf("File: %v is valid", file.Filename)
	default:
		logger.Errorf("Filename: %s, ext: %s", file.Filename, filepath.Ext(file.Filename))
		response.Error(c, response.JSONBody{Err: response.ErrWrongFileType, Msg: "Wrong file type, please upload a file of .crt or .key"}, http.StatusUnsupportedMediaType)
		return
	}

	var dstPath string
	filename := fmt.Sprintf("%s_%s", utils.RandStr(16), file.Filename)
	if config.GlobalConfig.Server.DevMode {
		dstPath = filepath.Join("./nginx", SSLCertDir, filename)
	} else {
		dstPath = filepath.Join(config.GlobalConfig.NgxResDir, SSLCertDir, filename)
	}
	if err = c.SaveUploadedFile(file, dstPath); err != nil {
		logger.Error(err)
		response.Error(c, response.JSONBody{Err: response.ErrInternalError, Msg: "Error occurred when saving file"}, http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"filename": filename})
}

func PostSSLCert(c *gin.Context) {
	var params postSSLCertRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
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
