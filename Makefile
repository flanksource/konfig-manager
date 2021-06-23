default: build
NAME:=konfig-manager

VERSION_TAG=$(VERSION)-$(shell date +"%Y%m%d%H%M%S")


# Image URL to use all building/pushing image targets
IMG ?= flanksource/konfig-manager:$(VERSION)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build: next
	go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o bin/konfig-manager

.PHONY: install
install: build
	cp bin/konfig-manager /usr/local/bin/

.PHONY: next
next:
	cd ui && npm ci && npm run export

.PHONY: linux
linux: next
	GOOS=linux GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_linux-arm64

.PHONY: darwin
darwin: next
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\""  -o .bin/$(NAME)_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_darwin-arm64

.PHONY: windows
windows: next
	GOOS=windows GOARCH=amd64 go build -o ./.bin/$(NAME).exe -ldflags "-X \"main.version=$(VERSION_TAG)\""  main.go

.PHONY: release
release: linux darwin windows compress

.PHONY: compress
compress:
	upx -5 ./.bin/*

.PHONY: docker
docker:
	docker build ./ -t $(NAME)

.PHONY: test
test:
	go test   ./test/... -test.v
	go test ./controllers/... -test.v -ginkgo.v
	go test   ./pkg/... -test.v

.PHONY: lint
lint: fmt vet
	golangci-lint run

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=konfig-manager webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

##@ Build

run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

## Apply the manifests for the operator
.PHONY: install-operator
install-operator:
	$(KUSTOMIZE) build config/base | kubectl apply -f -

## Delete the manifests for the operator
.PHONY: uninstall-operator
uninstall-operator:
	$(KUSTOMIZE) build config/base | kubectl delete -f -

install-crd: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall-crd: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -


CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)


KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

.PHONY: lint-markdown
lint-markdown: 
	npm install markdownlint-cli@0.27.1
	npx markdownlint '**/*.md' --ignore 'node_modules' --ignore 'ui/node_modules' -c .markdownlint.json

.PHONY: fix-markdown
fix-markdown: 
	npm install markdownlint-cli@0.27.1
	npx markdownlint '**/*.md' --fix --ignore 'node_modules' --ignore 'ui/node_modules' -c .markdownlint.json
