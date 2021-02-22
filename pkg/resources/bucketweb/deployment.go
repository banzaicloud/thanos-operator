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

package bucketweb

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
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

func (b *BucketWeb) deployment() (runtime.Object, reconciler.DesiredState, error) {
	meta := b.getMeta(b.getName())
	if b.ObjectStore.Spec.BucketWeb != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb
		deployment := &appsv1.Deployment{
			ObjectMeta: bucketWeb.MetaOverrides.Merge(meta),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: b.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: meta,
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  Name,
								Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(bucketWeb.HTTPAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "objectstore-secret",
										ReadOnly:  true,
										MountPath: "/etc/config/",
									},
								},
								ImagePullPolicy: corev1.PullIfNotPresent,
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "objectstore-secret",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: b.ObjectStore.Spec.Config.MountFrom.SecretKeyRef.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		var containerArgs = []string{
			"tools",
			"bucket",
			"web",
			"--log.level=info",
			"--http-address=" + bucketWeb.HTTPAddress,
			// TODO: get secret file path from secret mount
			"--objstore.config-file=/etc/config/" + b.ObjectStore.Spec.Config.MountFrom.SecretKeyRef.Key,
			"--refresh=" + strconv.Itoa(int(math.Floor(bucketWeb.Refresh.Duration.Seconds()))) + "s",
			"--timeout=" + strconv.Itoa(int(math.Floor(bucketWeb.Timeout.Duration.Seconds()))) + "s",
		}
		imageVersion := strings.Split(deployment.Spec.Template.Spec.Containers[0].Image, ":")[1]
		imgSem, _ := semver.NewVersion(imageVersion)
		breakingSem, _ := semver.NewConstraint("0.13.0")
		if breakingSem.Check(imgSem) {
			containerArgs = append([]string{"tools"}, containerArgs...)
		}
		if bucketWeb.Label != "" {
			containerArgs = append(containerArgs, "--label="+bucketWeb.Label)
		}
		deployment.Spec.Template.Spec.Containers[0].Args = containerArgs

		if bucketWeb.DeploymentOverrides != nil {
			if err := merge.Merge(deployment, bucketWeb.DeploymentOverrides); err != nil {
				return deployment, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to deployment base object")
			}
		}

		return deployment, reconciler.StatePresent, nil
	}

	return &appsv1.Deployment{
		ObjectMeta: meta,
	}, reconciler.StateAbsent, nil
}
