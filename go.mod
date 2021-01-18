module github.com/banzaicloud/thanos-operator

go 1.14

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.16.3-0.20210118211330-d8450836aad0
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.4
	github.com/go-logr/logr v0.2.1
	github.com/imdario/mergo v0.3.9
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.6.2
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
