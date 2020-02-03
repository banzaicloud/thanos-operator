package query

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *Query) deployment() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil {
		query := q.Thanos.Spec.Query.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: q.getMeta(q.getName()),
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
								Name:  "query",
								Image: fmt.Sprintf("%s:%s", query.Image.Repository, query.Image.Tag),
								Args: []string{
									"query",
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(query.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(query.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								Resources:       query.Resources,
								ImagePullPolicy: query.Image.PullPolicy,
								LivenessProbe:   q.GetCheck(resources.GetPort(query.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  q.GetCheck(resources.GetPort(query.HttpAddress), resources.ReadyCheckPath),
								VolumeMounts:    q.getVolumeMounts(),
							},
						},
						Volumes: q.getVolumes(),
					},
				},
			},
		}
		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = q.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)
		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}

func (q *Query) getVolumeMounts() []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}
	if q.Thanos.Spec.Query.GRPCClientCertificate != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "client-certificate",
			ReadOnly:  true,
			MountPath: serverCertMountPath,
		})
	}
	return volumeMounts
}

func (q *Query) getVolumes() []corev1.Volume {
	volumes := []corev1.Volume{}
	if q.Thanos.Spec.Query.GRPCClientCertificate != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "client-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: q.Thanos.Spec.Query.GRPCClientCertificate,
				},
			},
		})
	}
	return volumes
}
