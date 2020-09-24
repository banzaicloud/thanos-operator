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

package receiver

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *receiverInstance) service() (runtime.Object, reconciler.DesiredState, error) {
	if s.receiverGroup != nil {
		receiver := s.receiverGroup
		storeService := &corev1.Service{
			ObjectMeta: s.receiverGroup.MetaOverrides.Merge(s.getMeta()),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "grpc",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(receiver.GRPCAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "grpc",
						},
					},
					{
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(receiver.HTTPAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "http",
						},
					},
				},
				Selector:  s.getLabels(),
				ClusterIP: corev1.ClusterIPNone,
				Type:      corev1.ServiceTypeClusterIP,
			},
		}
		return storeService, reconciler.StatePresent, nil

	}
	delete := &corev1.Service{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
