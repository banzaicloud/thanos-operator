// Copyright 2021 Banzai Cloud
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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	ctrl "sigs.k8s.io/controller-runtime"
)

func NewReconciler(receiver *v1alpha1.Receiver, reconciler reconciler.ResourceReconciler) *Reconciler {
	return &Reconciler{
		Receiver:           receiver,
		ResourceReconciler: reconciler,
	}
}

type Reconciler struct {
	Receiver           *v1alpha1.Receiver
	ResourceReconciler reconciler.ResourceReconciler
}

func (r Reconciler) Reconcile() (*ctrl.Result, error) {
	resourceList, err := getResources(r.Receiver)
	if err != nil {
		return nil, err
	}
	return resources.Dispatch(r.ResourceReconciler, resourceList)
}

func getResources(receiver *v1alpha1.Receiver) (resourceList []resources.Resource, err error) {
	if receiver == nil {
		return
	}

	r := *New(receiver)
	resourceList = append(resourceList, (&receiverInstance{r, nil}).commonService)

	for _, group := range receiver.Spec.ReceiverGroups {
		if err = mergo.Merge(&group, v1alpha1.DefaultReceiverGroup); err != nil {
			return
		}
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).statefulset)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).hashring)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).service)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).serviceMonitor)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).ingressGRPC)
		resourceList = append(resourceList, (&receiverInstance{r, group.DeepCopy()}).ingressHTTP)
	}

	return
}
