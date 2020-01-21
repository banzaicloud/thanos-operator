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

package compactor

import (
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const Name = "compactor"

type Compactor struct {
	*resources.ObjectStoreReconciler
}

func New(reconciler *resources.ObjectStoreReconciler) *Compactor {
	return &Compactor{
		ObjectStoreReconciler: reconciler,
	}
}

func (c *Compactor) Reconcile() (*reconcile.Result, error) {
	if c.ObjectSore.Spec.Compactor != nil {
		err := mergo.Merge(c.ObjectSore.Spec.Compactor, v1alpha1.DefaultCompactor)
		if err != nil {
			return nil, err
		}
	}
	return c.ReconcileResources([]resources.Resource{
		//c.persistentVolumeClaim,
		c.deployment,
		c.service,
	})
}

func (c *Compactor) getLabels(name string) resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: name,
	}.Merge(c.GetCommonLabels())
	if c.ObjectSore.Spec.Compactor != nil {
		labels.Merge(c.ObjectSore.Spec.Compactor.Labels)
	}
	return labels
}

func (c *Compactor) getMeta(name string) metav1.ObjectMeta {
	meta := c.GetObjectMeta(name)
	meta.Labels = c.getLabels(name)
	if c.ObjectSore.Spec.Compactor != nil {
		meta.Annotations = c.ObjectSore.Spec.Compactor.Annotations
	}
	return meta
}
