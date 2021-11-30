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

package resources

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReceiverReconciler struct {
	*v1alpha1.Receiver
	reconciler.ResourceReconciler
}

func (t *ReceiverReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.Receiver.Name,
	}
}

func (t *ReceiverReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.Receiver.Name, name)
}

func (t *ReceiverReconciler) GetObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: t.Namespace,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: t.APIVersion,
				Kind:       t.Kind,
				Name:       t.Name,
				UID:        t.UID,
				Controller: utils.BoolPointer(true),
			},
		},
	}
}

func (t *ReceiverReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
	return Dispatch(t.ResourceReconciler, resourceList)
}

func NewReceiverReconciler(receiver *v1alpha1.Receiver, genericReconciler reconciler.ResourceReconciler) *ReceiverReconciler {
	return &ReceiverReconciler{
		Receiver:           receiver,
		ResourceReconciler: genericReconciler,
	}
}
