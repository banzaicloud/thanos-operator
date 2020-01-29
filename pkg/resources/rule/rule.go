package rule

import (
	"fmt"
	"strconv"

	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Rule struct {
	*resources.ThanosComponentReconciler
}

type ruleInstance struct {
	*Rule
	*v1alpha1.StoreEndpoint
}

func (r *ruleInstance) getName() string {
	return r.QualifiedName(fmt.Sprintf("%s-%s", r.StoreEndpoint.Name, v1alpha1.RuleName))
}

func (r *ruleInstance) getMeta() metav1.ObjectMeta {
	meta := r.GetNameMeta(r.getName(), "")
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.StoreEndpoint.APIVersion,
			Kind:       r.StoreEndpoint.Kind,
			Name:       r.StoreEndpoint.Name,
			UID:        r.StoreEndpoint.UID,
			Controller: util.BoolPointer(true),
		},
	}
	meta.Labels = r.Rule.getLabels()
	meta.Annotations = r.Thanos.Spec.Rule.Annotations
	return meta
}

func (r *ruleInstance) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc.cluster.local", r.getName(), r.StoreEndpoint.Namespace)
}

func New(reconciler *resources.ThanosComponentReconciler) *Rule {
	return &Rule{
		ThanosComponentReconciler: reconciler,
	}
}

func (r *Rule) resourceFactory() []resources.Resource {
	var resourceList []resources.Resource

	for _, endpoint := range r.StoreEndpoints {
		resourceList = append(resourceList, (&ruleInstance{r, endpoint.DeepCopy()}).statefulset)
		resourceList = append(resourceList, (&ruleInstance{r, endpoint.DeepCopy()}).service)
	}

	return resourceList
}
func (r *Rule) GetServiceURLS() []string {
	var urls []string
	for _, endpoint := range r.StoreEndpoints {
		urls = append(urls, (&ruleInstance{r, endpoint.DeepCopy()}).getSvc())
	}
	return urls
}

func (r *Rule) Reconcile() (*reconcile.Result, error) {
	if r.Thanos.Spec.Rule != nil {
		err := mergo.Merge(r.Thanos.Spec.Rule, v1alpha1.DefaultRule)
		if err != nil {
			return nil, err
		}
	}
	return r.ReconcileResources(r.resourceFactory())
}

func (r *Rule) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: v1alpha1.RuleName,
	}.Merge(
		r.GetCommonLabels(),
	)
	if r.Thanos.Spec.Rule != nil {
		labels.Merge(r.Thanos.Spec.Rule.Labels)
	}
	return labels
}

func (r *Rule) getQueryEndpoints() []string {
	var endpoints []string
	if r.Thanos.Spec.Query != nil {
		endpoints = append(endpoints, r.QualifiedName(v1alpha1.QueryName))
	}
	return endpoints
}

func (r *Rule) setArgs(args []string) []string {
	rule := r.Thanos.Spec.Rule.DeepCopy()
	args = append(args, resources.GetArgs(rule)...)

	for _, s := range r.getQueryEndpoints() {
		url := fmt.Sprintf("--query=%s.%s.svc.cluster.local:%s", s, r.Thanos.Namespace, strconv.Itoa(int(resources.GetPort(rule.HttpAddress))))
		args = append(args, url)
	}

	return args
}
