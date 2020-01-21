package rule

import (
	"fmt"
	"strconv"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const Name = "rule"

type Rule struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*resources.ThanosComponentReconciler
}

func (r *Rule) getLabels(name string) resources.Labels {
	return resources.Labels{
		resources.NameLabel: name,
	}.Merge(
		r.Thanos.Spec.Rule.Labels,
		r.GetCommonLabels(),
	)
}

func (r *Rule) getMeta(name string) metav1.ObjectMeta {
	meta := r.GetObjectMeta(name)
	meta.Labels = r.getLabels(name)
	meta.Annotations = r.Thanos.Spec.Rule.Annotations
	return meta
}

func (r *Rule) getQueryEndpoints() []string {
	var endpoints []string
	if r.Thanos.Spec.Query != nil {
		endpoints = append(endpoints, r.QualifiedName("query-service"))
	}
	return endpoints
}

func (r *Rule) setRuleArgs(args []string) []string {
	rule := r.Thanos.Spec.Rule.DeepCopy()
	args = append(args, resources.GetArgs(rule)...)
	if r.Thanos.Spec.ThanosDiscovery != nil {
		for _, s := range r.getQueryEndpoints() {
			url := fmt.Sprintf("--query=%s.%s.svc.cluster.local:%s", s, r.Thanos.Namespace, strconv.Itoa(int(resources.GetPort(rule.HttpAddress))))
			args = append(args, url)
		}
	}
	return args
}
