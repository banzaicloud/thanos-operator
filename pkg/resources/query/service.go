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

package query

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

func (q *Query) service() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil {
		queryService := &corev1.Service{
			ObjectMeta: q.Thanos.Spec.Query.MetaOverrides.Merge(q.getMeta(q.getName())),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "grpc",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(q.Thanos.Spec.Query.GRPCAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "grpc",
						},
					},
					{
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(q.Thanos.Spec.Query.HttpAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "http",
						},
					},
				},
				Selector: q.getLabels(),
				Type:     corev1.ServiceTypeClusterIP,
			},
		}

		if q.Thanos.Spec.Query.ServiceOverrides != nil {
			merged := &corev1.Service{}
			err := merge.Merge(queryService, q.Thanos.Spec.Query.ServiceOverrides, merged)
			if err != nil {
				return nil, nil, errors.WrapIf(err, "unable to merge overrides to base object")
			}
			queryService = merged
		}

		return queryService, common.ServiceUpdateHook(queryService), nil
	}
	delete := &corev1.Service{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}
