module github.com/banzaicloud/thanos-operator/pkg/sdk

go 1.13

require (
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0-20200107212637-8598193c9ebb
	github.com/banzaicloud/operator-tools v0.4.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v11.0.1-0.20190516230509-ae8359b20417+incompatible
	sigs.k8s.io/controller-runtime v0.4.0
)
