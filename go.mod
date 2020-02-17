module github.com/banzaicloud/thanos-operator

go 1.13

require (
	emperror.dev/errors v0.4.2
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/banzaicloud/operator-tools v0.6.1-0.20200217204318-13d8a21c6ed8
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.0
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.8
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
	sigs.k8s.io/controller-runtime v0.4.0
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
