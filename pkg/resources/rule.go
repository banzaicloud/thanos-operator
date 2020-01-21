package resources

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (t *ThanosComponentReconciler) getRuleLabels() Labels {
	return Labels{
		nameLabel: "rule",
	}.merge(
		t.Thanos.Spec.Rule.Labels,
		t.getCommonLabels(),
	)
}

func (t *ThanosComponentReconciler) getRuleMeta(name string) metav1.ObjectMeta {
	meta := t.getObjectMeta(name)
	meta.Labels = t.getRuleLabels()
	meta.Annotations = t.Thanos.Spec.Rule.Annotations
	return meta
}

func (t *ThanosComponentReconciler) getQueryEndpoints() []string {
	var endpoints []string
	if t.Thanos.Spec.Query != nil {
		endpoints = append(endpoints, t.qualifiedName("query-service"))
	}
	return endpoints
}

func (t *ThanosComponentReconciler) setRuleArgs(args []string) []string {
	rule := t.Thanos.Spec.Rule.DeepCopy()
	args = append(args, getArgs(rule)...)
	if t.Thanos.Spec.ThanosDiscovery != nil {
		for _, s := range t.getQueryEndpoints() {
			args = append(args, fmt.Sprintf("--query=%s.%s.svc.cluster.local:%s", s, t.Thanos.Namespace))
		}
	}
	return args
}

func (t *ThanosComponentReconciler) ruleStatefulSet() (runtime.Object, reconciler.DesiredState, error) {
	name := "rule-statefulset"
	if t.Thanos.Spec.Rule != nil {
		rule := t.Thanos.Spec.Rule.DeepCopy()
		statefulset := &appsv1.StatefulSet{
			ObjectMeta: t.getRuleMeta(name),
			Spec: appsv1.StatefulSetSpec{
				Replicas: util.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: t.getRuleLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: t.getRuleMeta(name),
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
										ContainerPort: GetPort(rule.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: GetPort(rule.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								Resources:       rule.Resources,
								LivenessProbe:   t.getCheck(GetPort(rule.HttpAddress), healthCheckPath),
								ReadinessProbe:  t.getCheck(GetPort(rule.HttpAddress), readyCheckPath),
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
		ObjectMeta: t.getObjectMeta(name),
	}
	return delete, reconciler.StateAbsent, nil
}
