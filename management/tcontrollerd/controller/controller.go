package controller

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/log"
	pb "chaitin.cn/patronus/safeline-2/management/tcontrollerd/proto/website"
)

var (
	logger = log.GetLogger("controller")
)

func Handle() error {
	logger.Infof("Connect mgt-webserver at %s", config.GlobalConfig.MgtWebserver)
	gRPCConn, err := grpc.Dial(config.GlobalConfig.MgtWebserver, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)
	if err != nil {
		logger.Errorf("Fail to dial: %v", err)
		return err
	}

	wsClient := pb.NewWebsiteClient(gRPCConn)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logger.Errorf("Fail to close: %v", err)
			return
		}
	}(gRPCConn)

	if err = websiteHandler(wsClient); err != nil {
		return err
	}

	return nil
}
