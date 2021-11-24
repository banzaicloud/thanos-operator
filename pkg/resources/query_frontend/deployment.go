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

package query_frontend

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *QueryFrontend) deployment() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.QueryFrontend != nil {
		queryFrontend := q.Thanos.Spec.QueryFrontend

		deployment := &appsv1.Deployment{
			ObjectMeta: queryFrontend.MetaOverrides.Merge(q.getMeta(q.getName())),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: q.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: q.getMeta(q.getName()),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "query-frontend",
								Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
								Args: []string{
									"query-frontend",
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(queryFrontend.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								ImagePullPolicy: corev1.PullIfNotPresent,
								LivenessProbe:   resources.GetProbe(resources.GetPort(queryFrontend.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  resources.GetProbe(resources.GetPort(queryFrontend.HttpAddress), resources.ReadyCheckPath),
							},
						},
					},
				},
			},
		}

		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = q.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)

		if queryFrontend.DeploymentOverrides != nil {
			if err := merge.Merge(deployment, queryFrontend.DeploymentOverrides); err != nil {
				return deployment, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}
