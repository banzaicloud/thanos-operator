module github.com/banzaicloud/thanos-operator

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/semver v1.5.0
	github.com/banzaicloud/operator-tools v0.25.2-0.20210929211649-76d50b1e79f6
	github.com/banzaicloud/thanos-operator/pkg/sdk v0.0.0
	github.com/go-logr/logr v0.4.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	github.com/spf13/cast v1.3.1
	go.uber.org/zap v1.18.1
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.5
)

replace github.com/banzaicloud/thanos-operator/pkg/sdk => ./pkg/sdk
