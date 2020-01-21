// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bucketweb

import (
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BucketWeb struct {
	*resources.ObjectStoreReconciler
}

const Name = "bucket"

func New(reconciler *resources.ObjectStoreReconciler) *BucketWeb {
	return &BucketWeb{
		ObjectStoreReconciler: reconciler,
	}
}

func (c *BucketWeb) Reconcile() (*reconcile.Result, error) {
	return c.ReconcileResources([]resources.Resource{
		//c.persistentVolumeClaim,
		c.deployment,
		c.service,
	})
}

func (c *BucketWeb) getLabels(name string) resources.Labels {
	return resources.Labels{
		resources.NameLabel: name,
	}.Merge(
		c.ObjectStoreReconciler.ObjectSore.Spec.Compactor.Labels,
		c.GetCommonLabels(),
	)
}

func (c *BucketWeb) getMeta(name string) metav1.ObjectMeta {
	meta := c.GetObjectMeta(name)
	meta.Labels = c.getLabels(name)
	meta.Annotations = c.ObjectStoreReconciler.ObjectSore.Spec.Compactor.Annotations
	return meta
}
