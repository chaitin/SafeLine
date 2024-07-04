package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/model"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/ngcmd"
	pb "chaitin.cn/patronus/safeline-2/management/tcontrollerd/proto/website"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/utils"
)

const (
	Ping                   = "ping"
	Pong                   = "pong"
	EventTypeWebsite       = "website"
	EventTypeDeleteWebsite = "deleteWebsite"
	EventTypeFullWebsite   = "fullWebsite"

	nginxConfigPath       = "/etc/nginx/sites-enabled/" // 修改的为 host /resources/nginx/ 中的文件
	nginxFilePrefix       = "IF_backend_"
	nginxBackupFilePrefix = "BAK_IF_backend_"
	nginxFileMode         = 0644

	HttpsScheme      = "https"
	DefaultHttpsPort = "443"
)

var pong = pb.Response{Type: Pong, Msg: nil, Err: false}

func generateNginxConfig(website *model.WebsiteConfig) (string, error) {
	// only ONE upstream supported in v1.0
	var upstreamAddr string
	scheme := "http"
	for _, upstream := range website.Upstreams {
		urlInfo, err := url.Parse(upstream)
		if err != nil {
			return "", err
		}
		if urlInfo.Scheme == HttpsScheme && urlInfo.Port() == "" {
			urlInfo.Host = urlInfo.Host + ":" + DefaultHttpsPort
		}
		upstreamAddr = fmt.Sprintf(upstreamAddrTpl, urlInfo.Host)
		if len(urlInfo.Scheme) > 0 {
			scheme = urlInfo.Scheme
		}
	}

	sslFlag := ""
	sslCertFilename := ""
	sslCertKeyFilename := ""
	if website.KeyFilename != "" && website.CertFilename != "" {
		sslFlag = " ssl" // with a blank ahead
		sslCertFilename = fmt.Sprintf(certTpl, website.CertFilename)
		sslCertKeyFilename = fmt.Sprintf(certKeyTpl, website.KeyFilename)
	}

	// only ONE server name supported in v1.0
	var serverName string
	var addrAnyProperties string
	for _, sn := range website.ServerNames {
		if sn == "*" || sn == "" {
			sn = "_"
			addrAnyProperties = addrAnyPropertiesTpl
		}
		serverName = fmt.Sprintf(serverNameTpl, sn)
	}

	// only ONE port supported in v1.0
	var serverListen string
	for _, port := range website.Ports {
		serverListen = fmt.Sprintf(serverListenTpl, port, sslFlag, addrAnyProperties)
	}

	upstreamName := fmt.Sprintf(backendTpl, website.Id)
	nginxConfig := strings.TrimSpace(fmt.Sprintf(nginxConfigTpl, upstreamName, upstreamAddr, serverListen, serverName, sslCertFilename, sslCertKeyFilename, scheme, upstreamName, website.Id))
	return nginxConfig, nil
}

func nginxTestAndReload() error {
	err := ngcmd.NginxConfTest()
	if err != nil {
		return err
	}

	err = ngcmd.NginxConfReload()
	if err != nil {
		return err
	}

	return nil
}

func generateFullConfigAndReload(msg []byte) error {
	var websites []model.WebsiteConfig
	if err := json.Unmarshal(msg, &websites); err != nil {
		return err
	}

	configFilename := make(map[string]struct{})
	for _, website := range websites {
		configFilename[fmt.Sprintf("%s%d", nginxFilePrefix, website.Id)] = struct{}{}
	}
	filepath.Walk(nginxConfigPath, func(path string, info fs.FileInfo, err error) error {
		_, ok := configFilename[info.Name()]
		if !ok && strings.HasPrefix(info.Name(), nginxFilePrefix) {
			if err := os.Remove(path); err != nil {
				// not return error only logged error in order to delete not exist website nginx conf
				logger.Warn(err)
			}
		}

		return nil
	})

	for _, website := range websites {
		if err := generateConfigAndReload(&website); err != nil {
			// trigger a full site push when tcd starts, ignore some site errors, and push as much as possible
			logger.Warn(err)
		}
	}

	return nil
}

func generateConfigAndReload(website *model.WebsiteConfig) error {
	nginxConfig, err := generateNginxConfig(website)
	if err != nil {
		return err
	}

	configPath := filepath.Join(nginxConfigPath, fmt.Sprintf("%s%d", nginxFilePrefix, website.Id))
	backupPath := filepath.Join(nginxConfigPath, fmt.Sprintf("%s%d", nginxBackupFilePrefix, website.Id))

	oldConfigExists, err := utils.FileExist(configPath)
	if err != nil {
		return err
	}

	if oldConfigExists {
		oldConfig, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}

		if string(oldConfig) == nginxConfig {
			logger.Info("No changes in the new website config, skip nginx -s reload")
			return nil
		}

		// tmp save old config to the backup path
		if err = utils.CopyFile(configPath, backupPath); err != nil {
			return err
		}
	}

	if err = utils.EnsureWriteFile(configPath, []byte(nginxConfig), nginxFileMode); err != nil {
		return err
	}

	if err = nginxTestAndReload(); err != nil {
		nginxError := err
		if err = os.Remove(configPath); err != nil {
			return err
		}
		if oldConfigExists {
			// new config err, restore old config
			if err = utils.CopyFile(backupPath, configPath); err != nil {
				return err
			}
			if err = os.Remove(backupPath); err != nil {
				return err
			}
		}
		return nginxError
	}

	if oldConfigExists {
		if err = os.Remove(backupPath); err != nil {
			return err
		}
	}

	return nil
}

func deleteConfigAndReload(config []byte) error {
	var websiteIds []uint
	if err := json.Unmarshal(config, &websiteIds); err != nil {
		return err
	}

	for _, id := range websiteIds {
		configPath := filepath.Join(nginxConfigPath, fmt.Sprintf("%s%d", nginxFilePrefix, id))
		exists, err := utils.FileExist(configPath)
		if err != nil {
			return err
		}
		if exists {
			if err := os.Remove(configPath); err != nil {
				return err
			}
		}
	}

	if err := nginxTestAndReload(); err != nil {
		return err
	}

	return nil
}

func sendResponse(stream pb.Website_SubscribeClient, eventType string, errMsg []byte) error {
	return stream.Send(&pb.Response{
		Type: eventType,
		Msg:  errMsg,
		Err:  len(errMsg) != 0,
	})
}

func websiteHandler(wc pb.WebsiteClient) error {
	stream, err := wc.Subscribe(context.Background())
	if err != nil {
		logger.Errorf("Subscribe failed: %v", err)
		return err
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// read done.
				logger.Infof("Recv EOF from webserver")
				return nil
			} else {
				logger.Errorf("Handle failed: %v", err)
				return err
			}
		}
		logger.Debugf("Got message Type %s, Msg: %s", event.GetType(), event.GetMsg())
		if event.Type == Ping {
			if err = stream.Send(&pong); err != nil {
				return err
			}
		} else if event.Type == EventTypeWebsite {
			logger.Infof("Update website with config: %s", event.GetMsg())
			var website model.WebsiteConfig
			if err := json.Unmarshal(event.GetMsg(), &website); err != nil {
				return err
			}

			if err = generateConfigAndReload(&website); err != nil {
				logger.Error(err)
				if err := sendResponse(stream, EventTypeDeleteWebsite, []byte(err.Error())); err != nil {
					return err
				}
			} else {
				if err = sendResponse(stream, EventTypeDeleteWebsite, nil); err != nil {
					return err
				}
			}
		} else if event.Type == EventTypeFullWebsite {
			if err = generateFullConfigAndReload(event.Msg); err != nil {
				logger.Error(err)
				sendErr := sendResponse(stream, EventTypeFullWebsite, []byte(err.Error()))
				if sendErr != nil {
					return sendErr
				}
			} else {
				if err = sendResponse(stream, EventTypeFullWebsite, nil); err != nil {
					return err
				}
			}
		} else if event.Type == EventTypeDeleteWebsite {
			if err = deleteConfigAndReload(event.Msg); err != nil {
				logger.Error(err)
				sendErr := sendResponse(stream, EventTypeDeleteWebsite, []byte(err.Error()))
				if sendErr != nil {
					return sendErr
				}
			} else {
				if err = sendResponse(stream, EventTypeDeleteWebsite, nil); err != nil {
					return err
				}
			}
		}
	}
}
