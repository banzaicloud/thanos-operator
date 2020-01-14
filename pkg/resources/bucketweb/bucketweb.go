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
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BucetWeb struct {
	client      client.Client
	namespace   string
	objectStore *v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func New(client client.Client, namespace string, objectStore *v1alpha1.ObjectStore, genericReconciler *reconciler.GenericResourceReconciler) *BucetWeb {
	return &BucetWeb{
		client:                    client,
		namespace:                 namespace,
		objectStore:               objectStore,
		GenericResourceReconciler: genericReconciler,
	}
}

func (b *BucetWeb) Reconcile() (*reconcile.Result, error) {
	for _, factory := range []resources.Resource{
		b.deployment,
		b.ingress,
		b.podDistributionBucket,
		b.service,
	} {
		o, state, err := factory()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", factory)
		}
		result, err := b.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapWithDetails(err,
				"failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (b *BucetWeb) objectMeta(name string, bucketWeb *v1alpha1.BaseObject) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   b.namespace,
		Labels:      bucketWeb.Labels,
		Annotations: bucketWeb.Annotations,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: b.objectStore.APIVersion,
				Kind:       b.objectStore.Kind,
				Name:       b.objectStore.Name,
				UID:        b.objectStore.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
}

func (b *BucetWeb) labels() map[string]string {
	const app = "bucketweb"

	return util.MergeLabels(b.objectStore.Spec.BucketWeb.Labels, map[string]string{
		"app.kubernetes.io/name":      app,
		"app.kubernetes.io/component": "bucket",
	}, generateLoggingRefLabels(b.objectStore.ObjectMeta.GetName()))
}

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func GetPort(address string) int32 {
	port, err := strconv.Atoi(strings.Split(address, ":")[1])
	if err != nil {
		return 0
	}
	return int32(port)
}
