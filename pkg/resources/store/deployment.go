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

func (s *Store) deployment() (runtime.Object, reconciler.DesiredState, error) {
	storeEndpoint, err := s.GetStoreEndpoint()
	if err != nil {
		return nil, nil, err
	}
	if s.Thanos.Spec.StoreGateway != nil {
		if storeEndpoint.Spec.Config.MountFrom == nil {
			return nil, nil, fmt.Errorf("missing config for StorageGateway %q", storeEndpoint.Name)
		}
		store := s.Thanos.Spec.StoreGateway.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: s.getMeta(v1alpha1.StoreName),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: s.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: s.getMeta(v1alpha1.StoreName),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  v1alpha1.StoreName,
								Image: fmt.Sprintf("%s:%s", store.Image.Repository, store.Image.Tag),
								Args: []string{
									"store",
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", storeEndpoint.Spec.Config.MountFrom.SecretKeyRef.Key),
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
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "objectstore-secret",
										ReadOnly:  true,
										MountPath: "/etc/config/",
									},
								},
								Resources:       store.Resources,
								ImagePullPolicy: store.Image.PullPolicy,
								LivenessProbe:   s.GetCheck(resources.GetPort(store.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  s.GetCheck(resources.GetPort(store.HttpAddress), resources.ReadyCheckPath),
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "objectstore-secret",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: storeEndpoint.Spec.Config.MountFrom.SecretKeyRef.Name,
									},
								},
							},
						},
					},
				},
			},
		}
		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = s.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)
		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: s.getMeta(v1alpha1.StoreName),
	}
	return delete, reconciler.StateAbsent, nil
}
