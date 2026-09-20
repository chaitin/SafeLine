package controller

import (
	"context"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/log"
	pb "chaitin.cn/patronus/safeline-2/management/tcontrollerd/proto/website"
)

var (
	logger = log.GetLogger("controller")
)

const (
	// TokenEnv holds the control channel token itself.
	TokenEnv = "TCD_MGT_TOKEN"

	// tokenMetadataKey is the metadata key that carries the token to the
	// management server. gRPC metadata keys have to be lower case.
	tokenMetadataKey = "mgt-token"
)

func Handle() error {
	logger.Infof("Connect mgt-webserver at %s", config.GlobalConfig.MgtWebserver)

	// The management server only accepts a subscription that proves it belongs
	// to this installation.
	token, err := controlToken()
	if err != nil {
		logger.Errorf("Fail to read the control channel token: %v", err)
		return err
	}

	dialOptions := []grpc.DialOption{}

	// The token below proves that this is tcontrollerd, and the certificate
	// proves that the peer is the management server of this installation:
	// reading the traffic of the deployment is not enough to take either side.
	creds, err := transportCreds()
	if err != nil {
		return err
	}
	if creds != nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	} else {
		// nil creds only when TCD_MGT_TLS is an explicit off value.
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	gRPCConn, err := grpc.Dial(config.GlobalConfig.MgtWebserver, dialOptions...)
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

	ctx := metadata.AppendToOutgoingContext(context.Background(), tokenMetadataKey, token)
	if err = websiteHandler(ctx, wsClient); err != nil {
		return err
	}

	return nil
}

// controlToken returns the shared secret of the control channel.
//
// The value is taken from the environment when it is set there, and from the
// token file otherwise. The management server creates that file in the
// directory both containers mount, so the standard deployment needs no
// configuration.
func controlToken() (string, error) {
	if token := strings.TrimSpace(os.Getenv(TokenEnv)); token != "" {
		return token, nil
	}

	content, err := os.ReadFile(config.GlobalConfig.MgtTokenFile)
	if err != nil {
		return "", errors.Wrapf(err, "read the control channel token from %s", config.GlobalConfig.MgtTokenFile)
	}

	token := strings.TrimSpace(string(content))
	if token == "" {
		return "", errors.New("the control channel token file " + config.GlobalConfig.MgtTokenFile + " is empty")
	}

	return token, nil
}
