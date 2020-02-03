package store

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *storeInstance) deployment() (runtime.Object, reconciler.DesiredState, error) {
	storeGateway := s.Store.Thanos.Spec.StoreGateway
	if storeGateway != nil {
		if s.StoreEndpoint.Spec.Config.MountFrom == nil {
			return nil, nil, fmt.Errorf("missing config for StorageGateway %q", s.StoreEndpoint.Name)
		}
		store := storeGateway.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: s.getMeta(),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: s.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: s.getMeta(),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  v1alpha1.StoreName,
								Image: fmt.Sprintf("%s:%s", store.Image.Repository, store.Image.Tag),
								Args: []string{
									"store",
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", s.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Key),
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(store.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(store.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								VolumeMounts:    s.getVolumeMounts(),
								Resources:       store.Resources,
								ImagePullPolicy: store.Image.PullPolicy,
								LivenessProbe:   s.GetCheck(resources.GetPort(store.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  s.GetCheck(resources.GetPort(store.HttpAddress), resources.ReadyCheckPath),
							},
						},
						Volumes: s.getVolumes(),
					},
				},
			},
		}

		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = s.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)
		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}

func (s *storeInstance) getVolumeMounts() []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "objectstore-secret",
			ReadOnly:  true,
			MountPath: "/etc/config/",
		},
	}
	if s.Thanos.Spec.StoreGateway.GRPCServerCertificate != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "server-certificate",
			ReadOnly:  true,
			MountPath: serverCertMountPath,
		})
	}
	return volumeMounts
}

func (s *storeInstance) getVolumes() []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "objectstore-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: s.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Name,
				},
			},
		},
	}
	if s.Thanos.Spec.StoreGateway.GRPCServerCertificate != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "server-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: s.Thanos.Spec.StoreGateway.GRPCServerCertificate,
				},
			},
		})
	}
	return volumes
}
