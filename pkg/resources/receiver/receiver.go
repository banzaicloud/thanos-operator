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

package receiver

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Receiver struct {
	*resources.ReceiverReconciler
}

type receiverInstance struct {
	*Receiver
	receiverGroup *v1alpha1.ReceiverGroup
}

func (r *receiverInstance) getName(suffix ...string) string {
	name := r.QualifiedName(v1alpha1.ReceiverName)
	if len(suffix) > 0 && suffix[0] != "" {
		name = name + "-" + suffix[0]
	}
	return name
}

func (r *receiverInstance) getVolumeMeta(name string) metav1.ObjectMeta {
	meta := r.GetNameMeta(name, "")
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.APIVersion,
			Kind:       r.Kind,
			Name:       r.Name,
			UID:        r.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r *receiverInstance) getMeta(suffix ...string) metav1.ObjectMeta {
	nameSuffix := ""
	if len(suffix) > 0 {
		nameSuffix = suffix[0]
	}
	meta := r.GetNameMeta(r.getName(nameSuffix), "")
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.APIVersion,
			Kind:       r.Kind,
			Name:       r.Name,
			UID:        r.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r *receiverInstance) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc.cluster.local", r.getName(), r.receiverGroup.Namespace)
}

func New(reconciler *resources.ReceiverReconciler) *Receiver {
	return &Receiver{
		reconciler,
	}
}

func (r *Receiver) resourceFactory() ([]resources.Resource, error) {
	var resourceList []resources.Resource

	for _, group := range r.Spec.ReceiverGroups {
		err := mergo.Merge(&group, v1alpha1.DefaultReceiverGroup)
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).statefulset)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).hashring)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).service)
		//resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).serviceMonitor)
		//resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).ingressHTTP)
		//resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).ingressGRPC)
	}

	return resourceList, nil
}

//func (r *Receiver) GetServiceURLS() []string {
//	var urls []string
//	for _, endpoint := range r.StoreEndpoints {
//		urls = append(urls, (&receiverInstance{r, endpoint.DeepCopy()}).getSvc())
//	}
//	return urls
//}

func (r *Receiver) Reconcile() (*reconcile.Result, error) {
	resources, err := r.resourceFactory()
	if err != nil {
		return nil, err
	}
	return r.ReconcileResources(resources)
}

func (r *receiverInstance) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: v1alpha1.ReceiverName,
	}.Merge(
		r.GetCommonLabels(),
	)
	return labels
}

func (r *receiverInstance) setArgs(args []string) []string {
	args = append(args, resources.GetArgs(r.receiverGroup)...)

	//Label

	// Local-endpoint

	return args
}
