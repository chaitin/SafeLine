package rpc

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	pb "chaitin.cn/patronus/safeline-2/management/webserver/proto/website"
)

const (
	WaitRspTimeout   = 30 * time.Second
	KeepaliveTime    = 5 * time.Second
	KeepaliveTimeout = 30 * time.Second

	Ping                   = "ping"
	Pong                   = "pong"
	EventTypeWebsite       = "website"
	EventTypeDeleteWebsite = "deleteWebsite"
	EventTypeFullWebsite   = "fullWebsite"
)

var logger = log.GetLogger("grpc")

// ServerOptions returns the options the control channel is served with.
//
// The channel carries the complete site configuration of the installation and
// republishes it on every change, so a subscription has to prove that it is
// tcontrollerd: the client certificate proves that the caller holds a
// credential of this installation, and the token is the second gate, which is
// what keeps a caller that reaches the port without one out.
func ServerOptions() ([]grpc.ServerOption, error) {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(subscribeAuthInterceptor),
	}

	creds, err := ServerCreds()
	if err != nil {
		return nil, err
	}
	opts = append(opts, creds)
	return opts, nil
}

func StartGRPCSever() error {
	lis, err := net.Listen("tcp", config.GlobalConfig.GPRC.ListenAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to listen")
	}

	if err := InitControlToken(); err != nil {
		// Not fatal: the console keeps working, and every subscription is
		// refused with an explanation until the token is available.
		logger.Errorf("Failed to init the control channel token: %s", err)
	}

	opts, err := ServerOptions()
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWebsiteServer(grpcServer, GetWebsiteServer())
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			logger.Fatalln("Failed to server")
		}
	}()
	return nil
}
