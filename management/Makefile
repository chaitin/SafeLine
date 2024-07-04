GO         = GO111MODULE=on go
#GO         = GO111MODULE=on GOOS=linux GOARCH=amd64 go
GOBUILD    = $(GO) build -mod=readonly
GOTEST     = $(GO) test -v -p 1 -coverprofile=coverage-management.out

STAMP    = $(shell date +%s)
GITHASH  = $(shell git rev-parse --short=8 HEAD)
GITTAG   = $(shell git describe --tags --abbrev=0)

BUILDFLAGS := -ldflags "-X main.buildstamp=$(STAMP) -X main.githash=$(GITHASH) -X main.version=$(GITTAG)"
pkgs        = ./...

all: build-all

.PHONY: build-all
build-all: proto build-webserver build-tcd

.PHONY: build-webserver
build-webserver:
	cd webserver && $(GOBUILD) $(BUILDFLAGS) -o ../build/webserver main.go

.PHONY: build-tcd
build-tcd:
	cd tcontrollerd && CGO_ENABLED=0 $(GOBUILD) $(BUILDFLAGS) -o ../build/tcontrollerd main.go

.PHONY: test
test:
	$(GOTEST) -failfast $(pkgs)

.PHONY: proto
proto:
	@./scripts/genproto.sh

.PHONY: lint
lint:
# 'go list' needs to be executed before staticcheck to prepopulate the modules cache.
# Otherwise staticcheck might fail randomly for some reason not yet explained.
	$(GO) list -e -compiled -test=true -export=false -deps=true -find=false -tags= -- ./... > /dev/null
	goimports -local chaitin.cn -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	golangci-lint version
	cd webserver && golangci-lint run -v --skip-dirs vendor --deadline 10m

.PHONY: clean
clean:
	rm -rf build

