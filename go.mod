module github.com/banzaicloud/thanos-operator

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.12.2
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.4
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.9
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
