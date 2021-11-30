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

type ObjectStoreReconciler struct {
	ObjectStore *v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func (t *ObjectStoreReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.ObjectStore.Name,
	}
}

func (t *ObjectStoreReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.ObjectStore.Name, name)
}

func (t *ObjectStoreReconciler) GetObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: t.ObjectStore.Namespace,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: t.ObjectStore.APIVersion,
				Kind:       t.ObjectStore.Kind,
				Name:       t.ObjectStore.Name,
				UID:        t.ObjectStore.UID,
				Controller: utils.BoolPointer(true),
			},
		},
	}
}

func (t *ObjectStoreReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
	return Dispatch(t.GenericResourceReconciler, resourceList)
}

func NewObjectStoreReconciler(objectStore *v1alpha1.ObjectStore, genericReconciler *reconciler.GenericResourceReconciler) *ObjectStoreReconciler {
	return &ObjectStoreReconciler{
		ObjectStore:               objectStore,
		GenericResourceReconciler: genericReconciler,
	}
}
