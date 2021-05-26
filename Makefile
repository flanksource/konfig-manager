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

.PHONY: release
release: release-darwin release-linux

.PHONY: release-linux
release-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o bin/$(NAME)_linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o bin/$(NAME)_linux-arm64

.PHONY: release-darwin
release-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION)\""  -o bin/$(NAME)_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION)\"" -o bin/$(NAME)_darwin-arm64

.PHONY: docker
docker:
	docker build ./ -t $(NAME)

.PHONY: e2e-tests
e2e-tests: build
	go run test/e2e.go

.PHONY: lint
lint: fmt vet
	golangci-lint run