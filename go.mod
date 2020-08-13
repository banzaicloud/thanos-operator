module github.com/banzaicloud/thanos-operator

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.10.5-0.20200812141632-f2d478750074
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.4
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.9
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
