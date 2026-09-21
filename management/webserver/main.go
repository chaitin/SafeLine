package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api"
	"chaitin.cn/patronus/safeline-2/management/webserver/cmd"
	"chaitin.cn/patronus/safeline-2/management/webserver/middleware"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/cron"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/rpc"
)

var (
	logger = log.GetLogger("main")

	version    = "undefined"
	githash    = "undefined"
	buildstamp = "undefined"
	goVersion  = "undefined"
)

func init() {
	// do something that do not raise error
}

func main() {
	log.SetLogFormatter()

	fs := flag.NewFlagSet("webserver", flag.ExitOnError)
	showVersion := fs.Bool("v", false, "show version")
	cfgFile := fs.String("c", constants.ConfigFilePath, "config file path")
	genCerts := fs.Bool("gen_certs", false, "generate certs")
	showFSL := fs.Bool("show_fsl", false, "show full selectors")
	push_fsl := fs.Bool("push_fsl", false, "compile and push fsl")
	fakeLogs := fs.Bool("fake_logs", false, "fake logs")
	resetUsername := fs.String("reset_user", "", "reset user")
	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatalln("Failed to parse args: ", err)
	}

	if *showVersion {
		i, _ := strconv.Atoi(buildstamp)
		t := time.Unix(int64(i), 0).Format("2006-01-02 15:04:05")
		fmt.Println("Version:    ", version)
		fmt.Println("Githash:    ", githash)
		fmt.Println("Build:      ", t)
		fmt.Println("Go version: ", goVersion)
		return
	}

	constants.Version = strings.TrimPrefix(version, "ce-")
	// init configs
	if err := config.InitConfigs(*cfgFile); err != nil {
		logger.Fatalln("Failed to init configs: ", err)
	}

	if err := log.InitLogger(); err != nil {
		logger.Fatalln("Failed to init db: ", err)
	}

	if *genCerts {
		if err := cmd.GenCerts(); err != nil {
			logger.Fatalln("Failed to generate certs: ", err)
		}
		return
	}

	logger.Info("Init database")
	if err := database.InitDB(); err != nil {
		logger.Fatalln("Failed to init db: ", err)
	}

	if *showFSL {
		if fullFSL, err := cmd.ShowFSL(); err != nil {
			logger.Fatalln("Failed to generate fsl: ", err)
		} else {
			logger.Info(strings.ReplaceAll(strings.ReplaceAll(fullFSL, ";", ";\n"), "CREATE", "\nCREATE"))
		}
		return
	}

	logger.Info("Init models")
	if err := model.InitModels(); err != nil {
		logger.Fatalln("Failed to init models: ", err)
	}

	if len(*resetUsername) > 0 {
		logger.Infoln("reset user:", *resetUsername)
		cmd.ResetUser(*resetUsername)
		logger.Infoln("success!")
		return
	}

	if *fakeLogs {
		logger.Infoln("faking logs...")
		cmd.FakeLogs()
		logger.Infoln("success!")
		return
	}

	if *push_fsl {
		logger.Infoln("push fsl...")
		if err := cmd.PushFSL(); err != nil {
			logger.Fatalln("Failed to generate fsl: ", err)
		}
		logger.Infoln("success!")
		return
	}

	logger.Info("Init FVM bytecode")
	fvm.InitFVMBytecode()

	if err := cron.StartCron(); err != nil {
		logger.Fatalln("Failed to start cron: ", err)
	}

	if err := rpc.StartGRPCSever(); err != nil {
		logger.Fatalln("Failed to start gRPC server: ", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// gin.ClientIP() returns the first value of X-Forwarded-For as soon as the
	// direct peer belongs to the trusted set, and the default trusted set is
	// every address. The management port is published on the host and reached
	// directly, so a client could otherwise pick the address it is throttled by
	// - the login throttle and the TFA bootstrap both key on it. No proxy is
	// trusted unless a deployment says so.
	if err := r.SetTrustedProxies(nil); err != nil {
		logger.Fatalln("Failed to disable trusted proxies: ", err)
	}
	r.Use(middleware.SecurityHeaders)
	// Nothing in front of this server caps the request body, so it is bounded
	// here (see middleware.BodyLimit).
	r.Use(middleware.BodyLimit(middleware.MaxRequestBodyBytes))

	var option model.Options
	database.GetDB().Where(&model.Options{Key: constants.SecretKey}).First(&option)
	store := cookie.NewStore([]byte(option.Value))
	r.Use(sessions.Sessions("session", store))

	created, tokenPath, err := api.EnsureBootstrapToken()
	if err != nil {
		logger.Fatalln("Failed to initialize the TFA bootstrap token: ", err)
	}
	if os.Getenv(api.BootstrapTokenEnv) != "" {
		logger.Infof("GET /api%s requires the %s header from %s.", api.OTPUrl, api.BootstrapTokenHeader, api.BootstrapTokenEnv)
	} else if created {
		logger.Infof("Wrote the TFA bootstrap token to %s (mode 0600). GET /api%s requires the %s header. The stock console must send it; until then bind with curl.", tokenPath, api.OTPUrl, api.BootstrapTokenHeader)
	} else {
		logger.Infof("GET /api%s requires the %s header; token loaded from %s.", api.OTPUrl, api.BootstrapTokenHeader, tokenPath)
	}
	if window := api.BootstrapWindow(); window <= 0 {
		logger.Warnf("%s is off: the account can be bound from the network until somebody logs in. A console that is left unbound stays claimable.",
			api.BootstrapWindowEnv)
	} else if openedAt, ok := model.BootstrapOpenedAt(); ok {
		logger.Infof("The TFA bootstrap of the account is open until %s (%s after it was opened at %s); run 'mgt -reset_user admin' to open it again.",
			openedAt.Add(window).Format(time.RFC3339), window, openedAt.Format(time.RFC3339))
	}

	publicRouters := r.Group("/api")
	publicRouters.Use(middleware.SameOrigin)
	publicRouters.POST(api.Login, api.PostLogin)
	publicRouters.POST(api.Logout, api.PostLogout)
	// First-boot QR still has no session, but requires X-Bootstrap-Token.
	publicRouters.GET(api.OTPUrl, api.GetOTPUrl)
	publicRouters.GET(api.Version, api.GetVersion)
	publicRouters.GET(api.UpgradeTips, api.GetUpgradeTips)
	// test use
	publicRouters.GET("/Ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	limitedRouters := r.Group("/api")
	limitedRouters.Use(middleware.SameOrigin)
	if envEnabled(os.Getenv("NO_AUTH")) {
		logger.Warn("Authentication is disabled because NO_AUTH is enabled")
	} else {
		limitedRouters.Use(middleware.AuthRequired)
	}
	if envEnabled(os.Getenv("READ_ONLY")) {
		logger.Warn("Write requests are rejected because READ_ONLY is enabled")
		limitedRouters.Use(middleware.ReadOnly)
	}

	limitedRouters.GET(api.User, api.GetUser)

	limitedRouters.POST(api.Behaviour, api.PostBehaviour)

	limitedRouters.GET(api.DetectLogList, api.GetDetectLogList)
	limitedRouters.GET(api.DetectLogDetail, api.GetDetectLogDetail)
	limitedRouters.POST(api.FalsePositives, api.PostFalsePositives)

	limitedRouters.POST(api.Website, api.PostWebsite)
	limitedRouters.PUT(api.Website, api.PutWebsite)
	limitedRouters.DELETE(api.Website, api.DeleteWebsite)
	limitedRouters.GET(api.Website, api.GetWebsite)

	limitedRouters.POST(api.UploadSSLCert, api.PostUploadSSLCert)
	limitedRouters.POST(api.SSLCert, api.PostSSLCert)

	limitedRouters.POST(api.PolicyRule, api.PostPolicyRule)
	limitedRouters.PUT(api.PolicyRule, api.PutPolicyRule)
	limitedRouters.DELETE(api.PolicyRule, api.DeletePolicyRule)
	limitedRouters.GET(api.PolicyRule, api.GetPolicyRule)
	limitedRouters.PUT(api.SwitchPolicyRule, api.PutSwitchPolicyRule)

	// 仪表盘接口
	limitedRouters.GET(api.DashboardCounts, api.GetDashboardCounts)
	limitedRouters.GET(api.DashboardSites, api.GetDashboardSites)
	limitedRouters.GET(api.DashboardQps, api.GetDashboardQps)
	limitedRouters.GET(api.DashboardRequests, api.GetDashboardRequests)
	limitedRouters.GET(api.DashboardIntercepts, api.GetDashboardIntercepts)

	limitedRouters.GET(api.PolicyGroupGlobal, api.GetPolicyGroupGlobal)
	limitedRouters.PUT(api.PolicyGroupGlobal, api.PutPolicyGroupGlobal)

	limitedRouters.GET(api.SrcIPConfig, api.GetSrcIPConfig)
	limitedRouters.PUT(api.SrcIPConfig, api.PutSrcIPConfig)

	logger.Info("Staring...")
	if err := r.Run(config.GlobalConfig.Server.ListenAddr); err != nil {
		logger.Fatalln("Error occurred when running web server: ", err)
	}
}

// envEnabled reports whether an environment variable holds an explicit truthy
// value. Switches such as NO_AUTH and READ_ONLY disable a protection, so only
// values that actually ask for it may enable them: an empty string (or any
// other unparsable value) must never do it.
func envEnabled(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
