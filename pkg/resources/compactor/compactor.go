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

package compactor

import (
	"github.com/banzaicloud/operator-tools/pkg/utils"
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
	if c.ObjectStore.Spec.Compactor != nil {
		err := mergo.Merge(c.ObjectStore.Spec.Compactor, v1alpha1.DefaultCompactor)
		if err != nil {
			return nil, err
		}
	}
	return c.ReconcileResources([]resources.Resource{
		//c.persistentVolumeClaim,
		c.deployment,
		c.service,
		c.serviceMonitor,
		c.persistentVolumeClaim,
	})
}

func (c *Compactor) getName() string {
	return c.QualifiedName(Name)
}

func (c *Compactor) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel: c.getName(),
		},
		c.GetCommonLabels(),
	)
}

func (c *Compactor) getMeta() metav1.ObjectMeta {
	meta := c.GetObjectMeta(c.getName())
	meta.Labels = c.getLabels()
	return meta
}
