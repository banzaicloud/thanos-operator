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
	"github.com/banzaicloud/thanos-operator/pkg/resources/common"
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
						TargetPort: intstr.IntOrString{IntVal: resources.GetPort(compactor.HTTPAddress)},
					},
				},
				Selector: c.getLabels(),
			},
		}

		if compactor.ServiceOverrides != nil {
			merged := &corev1.Service{}
			err := merge.Merge(service, compactor.ServiceOverrides, merged)
			if err != nil {
				return nil, nil, errors.WrapIf(err, "unable to merge overrides to service base")
			}
			service = merged
		}

		return service, common.ServiceUpdateHook(service), nil
	}

	return &corev1.Service{
		ObjectMeta: c.getMeta(),
	}, reconciler.StateAbsent, nil
}
