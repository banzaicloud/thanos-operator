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
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *storeInstance) service() (runtime.Object, reconciler.DesiredState, error) {
	if s.Thanos.Spec.StoreGateway != nil {
		store := s.Thanos.Spec.StoreGateway
		service := &corev1.Service{
			ObjectMeta: s.StoreEndpoint.Spec.MetaOverrides.Merge(s.getMeta()),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "grpc",
						Protocol:   corev1.ProtocolTCP,
						Port:       resources.GetPort(store.GRPCAddress),
						TargetPort: intstr.FromString("grpc"),
					},
					{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       resources.GetPort(store.HttpAddress),
						TargetPort: intstr.FromString("http"),
					},
				},
				Selector:  s.getLabels(),
				ClusterIP: corev1.ClusterIPNone,
				Type:      corev1.ServiceTypeClusterIP,
			},
		}
		if store.ServiceOverrides != nil {
			err := merge.Merge(service, store.ServiceOverrides)
			if err != nil {
				return service, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}
		return service, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
