package rule

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Rule) ruleStatefulSet() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil {
		rule := r.Thanos.Spec.Rule.DeepCopy()
		statefulset := &appsv1.StatefulSet{
			ObjectMeta: r.getMeta(Name),
			Spec: appsv1.StatefulSetSpec{
				Replicas: util.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: r.getLabels(Name),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: r.getMeta(Name),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "rule",
								Image: fmt.Sprintf("%s:%s", rule.Image.Repository, rule.Image.Tag),
								Args: []string{
									"rule",
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
								Resources:       rule.Resources,
								LivenessProbe:   r.GetCheck(resources.GetPort(rule.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  r.GetCheck(resources.GetPort(rule.HttpAddress), resources.ReadyCheckPath),
								ImagePullPolicy: rule.Image.PullPolicy,
							},
						},
					},
				},
			},
		}
		statefulset.Spec.Template.Spec.Containers[0].Args = t.setRuleArgs(statefulset.Spec.Template.Spec.Containers[0].Args)
		return statefulset, reconciler.StatePresent, nil
	}
	delete := &appsv1.StatefulSet{
		ObjectMeta: t.getObjectMeta(ruleName),
	}
	return delete, reconciler.StateAbsent, nil
}
