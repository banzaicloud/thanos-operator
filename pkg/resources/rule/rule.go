package rule

import (
	"fmt"
	"strconv"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const Name = "rule"

type Rule struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*resources.ThanosComponentReconciler
}

func New(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, reconciler *resources.ThanosComponentReconciler) *Rule {
	return &Rule{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		ThanosComponentReconciler: reconciler,
	}
}

func (r *Rule) Reconcile() (*reconcile.Result, error) {
	if r.Thanos.Spec.Rule != nil {
		err := mergo.Merge(r.Thanos.Spec.Rule, v1alpha1.DefaultRule)
		if err != nil {
			return nil, err
		}
	}
	return r.ReconcileResources(
		[]resources.Resource{
			r.statefulset,
			r.service,
		})
}

func (r *Rule) getLabels(name string) resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: name,
	}.Merge(
		r.GetCommonLabels(),
	)
	if r.Thanos.Spec.Rule != nil {
		labels.Merge(r.Thanos.Spec.Rule.Labels)
	}
	return labels
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

func (r *Rule) setArgs(args []string) []string {
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
