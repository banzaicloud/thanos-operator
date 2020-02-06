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

package sidecar

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(storeEndpoints *v1alpha1.StoreEndpointList, reconciler *resources.StoreEndpointComponentReconciler) *Sidecar {
	return &Sidecar{
		StoreEndpoints:                   storeEndpoints.Items,
		StoreEndpointComponentReconciler: reconciler,
	}
}

type Sidecar struct {
	StoreEndpoints []v1alpha1.StoreEndpoint
	*resources.StoreEndpointComponentReconciler
}

func (e *endpointService) getName() string {
	return fmt.Sprintf("%s-%s", e.Name, v1alpha1.SidecarName)
}

func (e *endpointService) getMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      e.getName(),
		Namespace: e.Namespace,
	}
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: e.APIVersion,
			Kind:       e.Kind,
			Name:       e.Name,
			UID:        e.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = e.Labels
	meta.Annotations = e.Annotations
	return meta
}

func (s *Sidecar) serviceFactory() []resources.Resource {
	var serviceList []resources.Resource

	for _, endpoint := range s.StoreEndpoints {
		serviceList = append(serviceList, (&endpointService{endpoint.DeepCopy()}).sidecarService)
		serviceList = append(serviceList, (&endpointService{endpoint.DeepCopy()}).ingressGRPC)
	}

	return serviceList
}

func (s *Sidecar) Reconcile() (*reconcile.Result, error) {
	return s.ReconcileResources(s.serviceFactory())
}
