# Web Server

web server for mgt-api

## Requirements

Go 1.18+

## Development

### Init protobuf

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

### Init fvm libs

```shell
# Due to the fvm c header files
mkdir -p management/webserver/submodule/fvm/
mkdir -p management/webserver/submodule/libct/

cd management/webserver/submodule/fvm/
# Download https://chaitin.cn/patronus/fvm/-/tags 1.8.21 release:release, https://chaitin.cn/patronus/fvm/-/jobs/6716645
unzip artifacts.zip
rm artifacts.zip

cd management/webserver/submodule/libct/
# Download https://chaitin.cn/patronus/libct/-/tags 1.1.1.0 release, https://chaitin.cn/patronus/libct/-/jobs/7229201
# rename
rm artifacts.zip

cd management/webserver/submodule/
# Download https://chaitin.cn/patronus/fusion-2/-/tags 5.3.9-r1 build:release, https://chaitin.cn/patronus/fusion-2/-/jobs/7326007
# rename
unzip artifacts.zip
mv artifacts/lib/libfusion.so libfvm.so
rm artifacts.zip
rm -r artifacts/
```

### Build

```shell
cd management/
docker run -it --rm -w="/mnt" --mount type=bind,source="$(pwd)",target=/mnt chaitin.cn/ci/golang:1.18 bash
cp webserver/submodule/libfvm.so /usr/lib/
make build-webserver
```