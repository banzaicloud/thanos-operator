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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type StoreEndpointComponentReconciler struct {
	StoreEndpoints []v1alpha1.StoreEndpoint
	*reconciler.GenericResourceReconciler
}

func (t *StoreEndpointComponentReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
	return Dispatch(t.GenericResourceReconciler, resourceList)
}

func NewStoreEndpointComponentReconciler(storeEndpoints *v1alpha1.StoreEndpointList, genericReconciler *reconciler.GenericResourceReconciler) *StoreEndpointComponentReconciler {
	return &StoreEndpointComponentReconciler{
		StoreEndpoints:            storeEndpoints.Items,
		GenericResourceReconciler: genericReconciler,
	}
}
