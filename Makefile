
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,maxDescLen=0"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOLANGCI_VERSION = 1.35.2
KUBEBUILDER_VERSION = 2.3.1
export KUBEBUILDER_ASSETS := $(PWD)/bin
LICENSEI_VERSION = 0.2.0

CONTROLLER_GEN_VERSION = v0.4.1
CONTROLLER_GEN = $(PWD)/bin/controller-gen

OS = $(shell uname | tr A-Z a-z)

all: manager

# Generate docs
.PHONY: docs
docs:
	go run cmd/docs.go

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) bin/golangci-lint
bin/golangci-lint-$(GOLANGCI_VERSION):
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run
	cd pkg/sdk && ../../bin/golangci-lint run

# Run tests
test: fmt vet genall lint bin/kubebuilder
	go test ./... -coverprofile cover.out
	cd pkg/sdk go test ./... -coverprofile cover-sdk.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate install
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
manifests: bin/controller-gen
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/thanos-operator/crds/
	cat config/rbac/role.yaml |sed 's|^  name: manager-role|  name: {{ include "thanos-operator.fullname" . }}|g' |sed '/creationTimestamp: null/d' | cat > charts/thanos-operator/templates/role.yaml

# Run go fmt against code
fmt:
	go fmt ./...
	cd pkg/sdk && go fmt ./...

# Run go vet against code
vet:
	go vet ./...
	cd pkg/sdk && go vet ./...

# Generate code
generate: bin/controller-gen
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

bin/controller-gen:
	@ if ! test -x bin/controller-gen; then \
		set -ex ;\
		CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
		cd $$CONTROLLER_GEN_TMP_DIR ;\
		go mod init tmp ;\
		GOBIN=$(PWD)/bin go get sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION} ;\
		rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	fi

# .PHONY: bin/kubebuilder_$(KUBEBUILDER_VERSION)
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

.PHONY: check
check: check-diff lint license-check test

install-minio:
	helm repo add minio https://helm.min.io/
	helm repo update
	helm upgrade --install minio minio/minio --set accessKey=myaccesskey,secretKey=mysecretkey,defaultBucket.enabled=true
	kubectl get secret thanos || kubectl create secret generic thanos --from-file=object-store.yaml=hack/object-store.yaml

install-prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	helm upgrade --install prometheus prometheus-community/kube-prometheus-stack -f hack/thanos-sidecar.yaml

install-thanos: install-minio install-prometheus
	kubectl apply -f config/samples/
