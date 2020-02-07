
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOLANGCI_VERSION = 1.19.1
KUBEBUILDER_VERSION = 2.2.0
export KUBEBUILDER_ASSETS := $(PWD)/bin
LICENSEI_VERSION = 0.2.0

OS = $(shell uname | tr A-Z a-z)

all: manager

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) bin/golangci-lint
bin/golangci-lint-$(GOLANGCI_VERSION):
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

# Run tests
test: generate fmt vet manifests bin/kubebuilder
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=$(IMG)
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths="./..."

genall: generate manifests docs
	go generate ./...
	cd pkg/sdk && go generate ./...

# Build the docker image
docker-build: test
	docker build . -t $(IMG)

# Push the docker image
docker-push:
	docker push $(IMG)

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: bin/kubebuilder_$(KUBEBUILDER_VERSION)
bin/kubebuilder_$(KUBEBUILDER_VERSION):
	@mkdir -p bin
	curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_amd64.tar.gz | tar xvz -C bin
	@ln -sf kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_amd64/bin bin/kubebuilder_$(KUBEBUILDER_VERSION)

bin/kubebuilder: bin/kubebuilder_$(KUBEBUILDER_VERSION)
	@ln -sf kubebuilder_$(KUBEBUILDER_VERSION)/kubebuilder bin/kubebuilder
	@ln -sf kubebuilder_$(KUBEBUILDER_VERSION)/kube-apiserver bin/kube-apiserver
	@ln -sf kubebuilder_$(KUBEBUILDER_VERSION)/etcd bin/etcd
	@ln -sf kubebuilder_$(KUBEBUILDER_VERSION)/kubectl bin/kubectl

check-diff:
	go mod tidy
	$(MAKE) genall
	git diff --exit-code

bin/licensei: bin/licensei-$(LICENSEI_VERSION)
	@ln -sf licensei-$(LICENSEI_VERSION) bin/licensei
bin/licensei-$(LICENSEI_VERSION):
	@mkdir -p bin
	curl -sfL https://git.io/licensei | bash -s v$(LICENSEI_VERSION)
	@mv bin/licensei $@

.PHONY: license-check
license-check: bin/licensei ## Run license check
	bin/licensei header
	bin/licensei check

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	bin/licensei cache
