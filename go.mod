module github.com/banzaicloud/thanos-operator

go 1.13

require (
	emperror.dev/errors v0.4.2
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0-20200107212637-8598193c9ebb
	github.com/banzaicloud/operator-tools v0.4.0
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.0
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.8
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	golang.org/x/tools v0.0.0-20200113040837-eac381796e91 // indirect
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v11.0.1-0.20190516230509-ae8359b20417+incompatible
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
)
