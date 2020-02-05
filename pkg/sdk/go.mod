module github.com/banzaicloud/thanos-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.4.2
	github.com/banzaicloud/operator-tools v0.4.2-0.20200205095409-b43b1ab90363
	github.com/cavaliercoder/grab v2.0.0+incompatible // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery v0.16.4
	sigs.k8s.io/controller-runtime v0.4.0
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1
