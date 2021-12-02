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
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *receiverInstance) service() (runtime.Object, reconciler.DesiredState, error) {
	if r.receiverGroup != nil {
		receiver := r.receiverGroup
		service := &corev1.Service{
			ObjectMeta: receiver.MetaOverrides.Merge(r.getMeta(receiver.Name)),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "grpc",
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(resources.GetPort(receiver.GRPCAddress)),
						TargetPort: intstr.FromString("grpc"),
					},
					{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(resources.GetPort(receiver.HTTPAddress)),
						TargetPort: intstr.FromString("http"),
					},
					{
						Name:       "remote-write",
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(resources.GetPort(receiver.RemoteWriteAddress)),
						TargetPort: intstr.FromString("remote-write"),
					},
				},
				Selector:  r.getLabels(),
				ClusterIP: corev1.ClusterIPNone,
				Type:      corev1.ServiceTypeClusterIP,
			},
		}
		if receiver.ServiceOverrides != nil {
			err := merge.Merge(service, receiver.ServiceOverrides)
			if err != nil {
				return service, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}
		return service, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: r.getMeta(r.receiverGroup.Name),
	}
	return delete, reconciler.StateAbsent, nil
}

func (r *receiverInstance) commonService() (runtime.Object, reconciler.DesiredState, error) {
	service := &corev1.Service{
		ObjectMeta: r.getMeta(),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       10907,
					TargetPort: intstr.FromString("grpc"),
				},
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       10909,
					TargetPort: intstr.FromString("http"),
				},
				{
					Name:       "remote-write",
					Protocol:   corev1.ProtocolTCP,
					Port:       10908,
					TargetPort: intstr.FromString("remote-write"),
				},
			},
			Selector: r.getLabels(),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	return service, reconciler.StatePresent, nil
}
