package rule

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Rule) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	storeEndpoint, err := r.GetStoreEndpoint()
	if err != nil {
		return nil, nil, err
	}
	if r.Thanos.Spec.Rule != nil {
		rule := r.Thanos.Spec.Rule.DeepCopy()
		statefulset := &appsv1.StatefulSet{
			ObjectMeta: r.getMeta(v1alpha1.RuleName),
			Spec: appsv1.StatefulSetSpec{
				Replicas: util.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: r.getLabels(v1alpha1.RuleName),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: r.getMeta(v1alpha1.RuleName),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "rule",
								Image: fmt.Sprintf("%s:%s", rule.Image.Repository, rule.Image.Tag),
								Args: []string{
									"rule",
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", storeEndpoint.Spec.Config.MountFrom.SecretKeyRef.Key),
								},
								WorkingDir: "",
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(rule.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(rule.GRPCAddress),
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
								Resources:       rule.Resources,
								LivenessProbe:   r.GetCheck(resources.GetPort(rule.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  r.GetCheck(resources.GetPort(rule.HttpAddress), resources.ReadyCheckPath),
								ImagePullPolicy: rule.Image.PullPolicy,
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
		statefulset.Spec.Template.Spec.Containers[0].Args = r.setArgs(statefulset.Spec.Template.Spec.Containers[0].Args)
		return statefulset, reconciler.StatePresent, nil
	}
	delete := &appsv1.StatefulSet{
		ObjectMeta: r.getMeta(v1alpha1.RuleName),
	}
	return delete, reconciler.StateAbsent, nil
}
