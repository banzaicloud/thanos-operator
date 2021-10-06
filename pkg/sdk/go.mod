module github.com/banzaicloud/thanos-operator/pkg/sdk

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/banzaicloud/operator-tools v0.25.4
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1
