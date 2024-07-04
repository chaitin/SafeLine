package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/controller"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/cron"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/ngcmd"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/log"
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

func handleLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			os.Exit(0)
		default:
			if err := controller.Handle(); err != nil {
				logger.Error("Error occurred when handling controller: ", err)
			}
			time.Sleep(time.Second * 5)
		}
	}
}

func main() {
	log.SetLogFormatter()

	fs := flag.NewFlagSet("tcontrollerd", flag.ExitOnError)
	showVersion := fs.Bool("v", false, "show version")
	nginxTest := fs.Bool("t", false, "nginx -t")
	nginxReload := fs.Bool("r", false, "nginx -s reload")
	cfgFile := fs.String("c", constants.ConfigFilePath, "config file path")
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

	// init configs
	if err := config.InitConfigs(*cfgFile); err != nil {
		logger.Fatalln("Failed to init configs: ", err)
	}

	if err := log.InitLogger(); err != nil {
		logger.Fatalln("Failed to init db: ", err)
	}

	if *nginxTest {
		if err := ngcmd.NginxConfTest(); err != nil {
			logger.Fatalln("Failed to test nginx conf: ", err)
		}
		return
	}

	if *nginxReload {
		if err := ngcmd.NginxConfReload(); err != nil {
			logger.Fatalln("Failed to reload nginx conf: ", err)
		}
		return
	}

	if err := cron.StartCron(); err != nil {
		logger.Fatalln("Failed to start cron: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	handleLoop(ctx)
}
