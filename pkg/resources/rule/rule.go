// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule

import (
	"fmt"
	"strconv"

	"github.com/banzaicloud/operator-tools/pkg/utils"
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

func (r *ruleInstance) getName(suffix ...string) string {
	name := r.QualifiedName(fmt.Sprintf("%s-%s", r.StoreEndpoint.Name, v1alpha1.RuleName))
	if len(suffix) > 0 && suffix[0] != "" {
		name = name + "-" + suffix[0]
	}
	return name
}

func (r *ruleInstance) getVolumeMeta(name string) metav1.ObjectMeta {
	meta := r.GetNameMeta(name)
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.StoreEndpoint.APIVersion,
			Kind:       r.StoreEndpoint.Kind,
			Name:       r.StoreEndpoint.Name,
			UID:        r.StoreEndpoint.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r *ruleInstance) getMeta(suffix ...string) metav1.ObjectMeta {
	nameSuffix := ""
	if len(suffix) > 0 {
		nameSuffix = suffix[0]
	}
	meta := r.GetNameMeta(r.getName(nameSuffix))
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.StoreEndpoint.APIVersion,
			Kind:       r.StoreEndpoint.Kind,
			Name:       r.StoreEndpoint.Name,
			UID:        r.StoreEndpoint.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r *ruleInstance) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc", r.getName(), r.StoreEndpoint.Namespace)
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
		resourceList = append(resourceList, (&ruleInstance{r, endpoint.DeepCopy()}).serviceMonitor)
		resourceList = append(resourceList, (&ruleInstance{r, endpoint.DeepCopy()}).ingressHTTP)
		resourceList = append(resourceList, (&ruleInstance{r, endpoint.DeepCopy()}).ingressGRPC)
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

func (r *ruleInstance) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel:     v1alpha1.RuleName,
			resources.StoreEndpoint: r.StoreEndpoint.Name,
		},
		r.GetCommonLabels(),
	)
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
		url := fmt.Sprintf("--query=%s.%s.svc:%s", s, r.Thanos.Namespace, strconv.Itoa(int(resources.GetPort(rule.HttpAddress))))
		args = append(args, url)
	}

	for key, value := range r.Thanos.Spec.Rule.Labels {
		args = append(args, fmt.Sprintf("--label %s=%s", key, value))
	}

	return args
}
