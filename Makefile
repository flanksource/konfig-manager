default: build
NAME:=konfig-manager

ifeq ($(VERSION),)
VERSION=$(shell git describe --tags  --long)-$(shell date +"%Y%m%d%H%M%S")
endif

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -ldflags "-X \"main.version=$(VERSION)\"" -o bin/konfig-manager


.PHONY: install
install: build
	cp bin/konfig-manager /usr/local/bin/

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o .bin/$(NAME)_linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o .bin/$(NAME)_linux-arm64

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION)\""  -o .bin/$(NAME)_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o .bin/$(NAME)_darwin-arm64

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build -o ./.bin/$(NAME).exe -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: release
release: linux darwin windows compress

.PHONY: compress
compress:
	# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
	which upx 2>&1 >  /dev/null  || (sudo apt-get update && sudo apt-get install -y xz-utils && wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; tar xf upx.tar.xz; mv upx-3.96-amd64_linux/upx /usr/bin )
	upx -5 ./.bin/$(NAME)_linux-amd64 ./.bin/$(NAME)_linux-arm64 ./.bin/$(NAME)_osx-amd64 ./.bin/$(NAME)_osx-arm64 ./.bin/$(NAME).exe

.PHONY: docker
docker:
	docker build ./ -t $(NAME)

.PHONY: test
test:
	go test   ./test/... -test.v

.PHONY: lint
lint: fmt vet
	golangci-lint run
