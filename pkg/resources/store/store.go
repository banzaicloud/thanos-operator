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

package store

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const serverCertMountPath = "/etc/tls/server"

func New(reconciler *resources.ThanosComponentReconciler) *Store {
	return &Store{
		ThanosComponentReconciler: reconciler,
	}
}

type storeInstance struct {
	*Store
	*v1alpha1.StoreEndpoint
}

func (s *storeInstance) getName() string {
	return s.QualifiedName(fmt.Sprintf("%s-%s", s.StoreEndpoint.Name, v1alpha1.StoreName))
}

func (s *storeInstance) getMeta() metav1.ObjectMeta {
	meta := s.GetNameMeta(s.getName())
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: s.StoreEndpoint.APIVersion,
			Kind:       s.StoreEndpoint.Kind,
			Name:       s.StoreEndpoint.Name,
			UID:        s.StoreEndpoint.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = s.getLabels()
	return meta
}

func (s *storeInstance) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc", s.getName(), s.StoreEndpoint.Namespace)
}

func (s *Store) resourceFactory() []resources.Resource {
	var resourceList []resources.Resource

	for _, endpoint := range s.StoreEndpoints {
		resourceList = append(resourceList, (&storeInstance{s, endpoint.DeepCopy()}).deployment)
		resourceList = append(resourceList, (&storeInstance{s, endpoint.DeepCopy()}).service)
		resourceList = append(resourceList, (&storeInstance{s, endpoint.DeepCopy()}).serviceMonitor)
		resourceList = append(resourceList, (&storeInstance{s, endpoint.DeepCopy()}).ingressGRPC)
	}

	return resourceList
}

func (s *Store) GetServiceURLS() []string {
	var urls []string
	for _, endpoint := range s.StoreEndpoints {
		urls = append(urls, (&storeInstance{s, endpoint.DeepCopy()}).getSvc())
	}
	return urls
}

func (s *Store) Reconcile() (*reconcile.Result, error) {
	if s.Thanos.Spec.StoreGateway != nil {
		err := mergo.Merge(s.Thanos.Spec.StoreGateway, v1alpha1.DefaultStoreGateway)
		if err != nil {
			return nil, err
		}
	}
	return s.ReconcileResources(s.resourceFactory())
}

type Store struct {
	*resources.ThanosComponentReconciler
}

func (s *storeInstance) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel:     v1alpha1.StoreName,
			resources.StoreEndpoint: s.StoreEndpoint.Name,
		},
		s.GetCommonLabels(),
	)
}

func (s *Store) setArgs(args []string) []string {
	store := s.Thanos.Spec.StoreGateway.DeepCopy()
	args = append(args, resources.GetArgs(store)...)
	if s.Thanos.Spec.StoreGateway.GRPCServerCertificate != "" {
		args = append(args, fmt.Sprintf("--grpc-server-tls-cert=%s/%s", serverCertMountPath, "tls.crt"))
		args = append(args, fmt.Sprintf("--grpc-server-tls-key=%s/%s", serverCertMountPath, "tls.key"))
		args = append(args, fmt.Sprintf("--grpc-server-tls-client-ca=%s/%s", serverCertMountPath, "ca.crt"))
	}
	return args
}
