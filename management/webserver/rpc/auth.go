package rpc

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	// TokenEnv holds the shared secret of the control channel between
	// tcontrollerd and this server.
	TokenEnv = "MGT_GRPC_TOKEN"

	// TokenFileEnv names the file the secret is read from when TokenEnv is
	// empty.
	TokenFileEnv = "MGT_GRPC_TOKEN_FILE"

	// DefaultTokenFile is where the secret is read from (and created) by
	// default. The standard deployment mounts the same host directory at
	// /app/sock in the management container and in the tengine container,
	// which is where tcontrollerd runs, so both sides find the same value
	// without any configuration.
	DefaultTokenFile = "/app/sock/mgt_grpc_token"

	// tokenMetadataKey is the gRPC metadata key that carries the secret. gRPC
	// metadata keys have to be lower case.
	tokenMetadataKey = "mgt-token"

	// tokenBytes is the size of a generated secret.
	tokenBytes = 32
)

// tokenFileMode keeps the secret readable for the owner only.
const tokenFileMode os.FileMode = 0600

// controlToken is the secret that Subscribe requires. It is empty when the
// server could not obtain one, in which case every subscription is refused.
var controlToken string

// InitControlToken loads the token of the control channel, creating it when it
// does not exist yet.
//
// A failure here does not stop the server: it disables the control channel,
// which is reported by Publish as an error, and leaves the console usable so
// that an operator can fix the deployment.
func InitControlToken() error {
	if token := strings.TrimSpace(os.Getenv(TokenEnv)); token != "" {
		controlToken = token
		return nil
	}

	token, err := loadOrCreateToken(TokenFileFromEnv())
	if err != nil {
		return err
	}

	controlToken = token
	return nil
}

// TokenFileFromEnv returns the path of the token file.
func TokenFileFromEnv() string {
	if path := strings.TrimSpace(os.Getenv(TokenFileEnv)); path != "" {
		return path
	}

	return DefaultTokenFile
}

// loadOrCreateToken reads the control channel token from path, creating it with
// a fresh random value when the file does not exist yet.
func loadOrCreateToken(path string) (string, error) {
	if token := readToken(path); token != "" {
		return token, nil
	}

	token := make([]byte, tokenBytes)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	value := hex.EncodeToString(token)

	// 0750: the directory holds the secret of the control channel, so it is not
	// left traversable by other users. It is a shared mount point whose mode is
	// normally set by the deployment, so this only applies when this call is the
	// one that creates it.
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", err
	}

	// O_EXCL: another process may have created the file between the read above
	// and this call, in which case its value wins.
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, tokenFileMode)
	if err != nil {
		if os.IsExist(err) {
			return readToken(path), nil
		}
		return "", err
	}

	if _, err = file.WriteString(value + "\n"); err != nil {
		_ = file.Close()
		return "", err
	}

	if err = file.Close(); err != nil {
		return "", err
	}

	return value, nil
}

// readToken returns the token stored in path, or an empty string when it is not
// available.
func readToken(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(content))
}

// subscribeAuthInterceptor refuses a Subscribe call that does not prove that it
// knows the control channel token.
//
// Without it, any client that can reach the gRPC port becomes the subscriber:
// it receives the full site configuration of the installation and silently
// takes the place of tcontrollerd, which then stops receiving its updates.
func subscribeAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if controlToken == "" {
		logger.Error("The control channel token is not configured, refusing subscriptions. Set " + TokenEnv + " or make " + TokenFileFromEnv() + " writable.")
		return status.Error(codes.Unauthenticated, "the control channel token is not configured on the server")
	}

	if subtle.ConstantTimeCompare([]byte(incomingToken(ss.Context())), []byte(controlToken)) != 1 {
		logger.Warnf("Rejected a subscription from %s: invalid control channel token", peerAddress(ss.Context()))
		return status.Error(codes.Unauthenticated, "a valid control channel token is required")
	}

	return handler(srv, ss)
}

// incomingToken returns the token that the client sent in the metadata.
func incomingToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(tokenMetadataKey)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// peerAddress returns the address of the client for the log.
func peerAddress(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		return p.Addr.String()
	}

	return "unknown"
}
