# TControllerD

Tengine Controller Daemon (abbr. TCD) will be running in the Tengine container, 
designed to be in place with minion on the host machine.

## Requirements

Go 1.18+

## Development

### init protobuf

```shell
# Refer: https://grpc.io/docs/languages/go/quickstart/
# 1. Install protoc
# https://grpc.io/docs/protoc-installation/

# 2. Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# 3. Update your PATH so that the protoc compiler can find the plugins
# export PATH="$PATH:$(go env GOPATH)/bin"

# 4. Generate proto go code
# cd /path/to/management
./scripts/genproto.sh
```