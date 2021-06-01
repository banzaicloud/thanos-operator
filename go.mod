module github.com/banzaicloud/thanos-operator

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.23.0
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.4
	github.com/go-logr/logr v0.4.0
	github.com/imdario/mergo v0.3.10
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	k8s.io/api v0.20.7
	k8s.io/apiextensions-apiserver v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk
