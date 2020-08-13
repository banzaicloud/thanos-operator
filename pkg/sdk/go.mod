module github.com/banzaicloud/thanos-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.4.2
	github.com/banzaicloud/operator-tools v0.10.5-0.20200812141632-f2d478750074
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/structured-merge-diff v1.0.1 // indirect
	sigs.k8s.io/testing_frameworks v0.1.2 // indirect
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1
