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
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Compactor) service() (runtime.Object, reconciler.DesiredState, error) {
	if c.ObjectStore.Spec.Compactor != nil {
		compactor := c.ObjectStore.Spec.Compactor
		service := &corev1.Service{
			ObjectMeta: compactor.MetaOverrides.Merge(c.getMeta()),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http",
						Port:       resources.GetPort(compactor.HTTPAddress),
						TargetPort: intstr.FromInt(int(resources.GetPort(compactor.HTTPAddress))),
					},
				},
				Selector: c.getLabels(),
			},
		}

		if compactor.ServiceOverrides != nil {
			err := merge.Merge(service, compactor.ServiceOverrides)
			if err != nil {
				return service, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to service base")
			}
		}

		return service, reconciler.DesiredStateHook(func(current runtime.Object) error {
			if s, ok := current.(*corev1.Service); ok {
				service.Spec.ClusterIP = s.Spec.ClusterIP
			} else {
				return errors.Errorf("failed to cast service object %+v", current)
			}
			return nil
		}), nil
	}

	return &corev1.Service{
		ObjectMeta: c.getMeta(),
	}, reconciler.StateAbsent, nil
}
