module github.com/banzaicloud/thanos-operator

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.10.3
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.4
	github.com/containerd/containerd v1.3.6
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.9
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
