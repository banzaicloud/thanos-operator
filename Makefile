
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,maxDescLen=0"

OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)

BIN := ${PWD}/bin

CONTROLLER_GEN := ${BIN}/controller-gen
CONTROLLER_GEN_VERSION = v0.5.0

ENVTEST_BIN_DIR := ${BIN}/envtest
ENVTEST_K8S_VERSION := 1.21.4
ENVTEST_BINARY_ASSETS := ${ENVTEST_BIN_DIR}/bin

GOLANGCI_LINT := ${BIN}/golangci-lint
GOLANGCI_LINT_VERSION := v1.42.1

KUBEBUILDER := ${BIN}/kubebuilder
KUBEBUILDER_VERSION = v3.1.0

LICENSEI := ${BIN}/licensei
LICENSEI_VERSION = v0.4.0

SETUP_ENVTEST := ${BIN}/setup-envtest

## =============
## ==  Rules  ==
## =============

all: manager

.PHONY: check
check: check-diff lint license-check test

.PHONY: check-circle
check-circle: check-diff lint test

.PHONY: check-diff
check-diff: genall tidy
	git diff --exit-code

.PHONY: deploy
deploy: manifests ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	cd config/manager && kustomize edit set image controller=$(IMG)
	kustomize build config/default | kubectl apply -f -

.PHONY: docker-build
docker-build: test ## Build the docker image
	docker build . -t $(IMG)

.PHONY: docker-push
docker-push: ## Push the docker image
	docker push $(IMG)

.PHONY: docs
docs: ## Generate docs
	go run cmd/docs.go

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...
	cd pkg/sdk && go fmt ./...

.PHONY: genall
genall: generate manifests docs tidy
	go generate ./...
	cd pkg/sdk && go generate ./...

.PHONY: generate
generate: ${CONTROLLER_GEN} ## Generate code
	cd pkg/sdk && ${CONTROLLER_GEN} object:headerFile=./../../hack/boilerplate.go.txt paths="./..."

.PHONY: help
help: ## Show this help message
	@grep -h -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' | sort

.PHONY: install
install: manifests ## Install CRDs into a cluster
	kustomize build config/crd | kubectl apply -f -

.PHONY: install-minio
install-minio:
	helm repo add minio https://helm.min.io/
	helm repo update
	helm upgrade --install minio minio/minio --set accessKey=myaccesskey,secretKey=mysecretkey,defaultBucket.enabled=true,resources.requests.memory=256Mi,persistence.size=5Gi
	kubectl get secret thanos || kubectl create secret generic thanos --from-file=object-store.yaml=hack/object-store.yaml

.PHONY: install-prometheus
install-prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	helm upgrade --install prometheus prometheus-community/kube-prometheus-stack -f hack/thanos-remote-write.yaml

.PHONY: install-thanos
install-thanos: install-minio install-prometheus
	kubectl apply -f config/samples/

.PHONY: license-cache
license-cache: ${LICENSEI} ## Generate license cache
	${LICENSEI} cache

.PHONY: license-check
license-check: ${LICENSEI} .licensei.cache ## Check licenses
	${LICENSEI} header
	${LICENSEI} check

.PHONY: lint
lint: ${GOLANGCI_LINT} ## Run linter
	${GOLANGCI_LINT} run ${LINTER_FLAGS}
	cd pkg/sdk && ${GOLANGCI_LINT} run ${LINTER_FLAGS}

.PHONY: list-all
list-all: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: manager
manager: fmt generate lint ## Build manager binary
	go build -o bin/manager main.go

.PHONY: manifests
manifests: ${CONTROLLER_GEN} # Generate manifests e.g. CRD, RBAC etc.
	cd pkg/sdk && ${CONTROLLER_GEN} $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	${CONTROLLER_GEN} $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/thanos-operator/crds/
	cat config/rbac/role.yaml |sed 's|^  name: manager-role|  name: {{ include "thanos-operator.fullname" . }}|g' |sed '/creationTimestamp: null/d' | cat > charts/thanos-operator/templates/role.yaml

.PHONY: run
run: generate install ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go

.PHONY: test
test: fmt genall lint ${ENVTEST_BINARY_ASSETS} ${KUBEBUILDER} ## Run tests
	ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./... -coverprofile cover.out
	cd pkg/sdk && go test ./... -coverprofile cover-sdk.out

tidy: ## Run `go mod tidy` against all modules
	find . -iname "go.mod" | xargs -L1 sh -c 'cd $$(dirname $$0); go mod tidy'

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from a cluster
	kustomize build config/crd | kubectl delete -f -

## =========================
## ==  Tool dependencies  ==
## =========================

${CONTROLLER_GEN}: ${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}: IMPORT_PATH := sigs.k8s.io/controller-tools/cmd/controller-gen
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}: VERSION := ${CONTROLLER_GEN_VERSION}
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}: | ${BIN}
	${go_install_binary}

${ENVTEST_BINARY_ASSETS}: ${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}
	ln -sf $(notdir $<) $@

${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}: | ${SETUP_ENVTEST} ${ENVTEST_BIN_DIR}
	ln -sf $$(${SETUP_ENVTEST} --bin-dir ${ENVTEST_BIN_DIR} use ${ENVTEST_K8S_VERSION} -p path) $@

${GOLANGCI_LINT}: ${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}: IMPORT_PATH := github.com/golangci/golangci-lint/cmd/golangci-lint
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}: VERSION := ${GOLANGCI_LINT_VERSION}
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}: | ${BIN}
	${go_install_binary}

${KUBEBUILDER}: ${KUBEBUILDER}_$(KUBEBUILDER_VERSION) | ${BIN}
	ln -sf $(notdir $<) $@

${KUBEBUILDER}_$(KUBEBUILDER_VERSION): | ${BIN}
	curl -sL https://github.com/kubernetes-sigs/kubebuilder/releases/download/${KUBEBUILDER_VERSION}/kubebuilder_${OS}_${ARCH} -o $@
	chmod +x $@

${LICENSEI}: ${LICENSEI}_${LICENSEI_VERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${LICENSEI}_${LICENSEI_VERSION}: IMPORT_PATH := github.com/goph/licensei/cmd/licensei
${LICENSEI}_${LICENSEI_VERSION}: VERSION := ${LICENSEI_VERSION}
${LICENSEI}_${LICENSEI_VERSION}: | ${BIN}
	${go_install_binary}

.licensei.cache: ${LICENSEI}
ifndef GITHUB_TOKEN
	@>&2 echo "WARNING: building licensei cache without Github token, rate limiting might occur."
	@>&2 echo "(Hint: If too many licenses are missing, try specifying a Github token via the environment variable GITHUB_TOKEN.)"
endif
	${LICENSEI} cache

${SETUP_ENVTEST}: IMPORT_PATH := sigs.k8s.io/controller-runtime/tools/setup-envtest
${SETUP_ENVTEST}: VERSION := latest
${SETUP_ENVTEST}: | ${BIN}
	GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}

${ENVTEST_BIN_DIR}: | ${BIN}
	mkdir -p $@

${BIN}:
	mkdir -p $@

define go_install_binary
find ${BIN} -name '$(notdir ${IMPORT_PATH})_*' -exec rm {} +
GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}
mv ${BIN}/$(notdir ${IMPORT_PATH}) ${BIN}/$(notdir ${IMPORT_PATH})_${VERSION}
endef
