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
package bucketweb

import (
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
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
	if c.ObjectStore.Spec.BucketWeb != nil {
		err := mergo.Merge(c.ObjectStore.Spec.BucketWeb, v1alpha1.DefaultBucketWeb)
		if err != nil {
			return nil, err
		}
	}
	return c.ReconcileResources([]resources.Resource{
		c.deployment,
		c.service,
		c.ingress,
	})
}

func (b *BucketWeb) getName() string {
	return b.QualifiedName(Name)
}

func (b *BucketWeb) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel: b.getName(),
		},
		b.GetCommonLabels(),
	)
}

func (b *BucketWeb) getMeta(name string) metav1.ObjectMeta {
	meta := b.GetObjectMeta(name)
	meta.Labels = b.getLabels()
	return meta
}
