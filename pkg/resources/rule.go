package resources

import (
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
		t.Thanos.Spec.Query.Labels,
		t.getCommonLabels(),
	)
}

func (t *ThanosComponentReconciler) getRuleMeta(name string) metav1.ObjectMeta {
	meta := t.getObjectMeta(name)
	meta.Labels = t.getRuleLabels()
	meta.Annotations = t.Thanos.Spec.Rule.Annotations
	return meta
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
					MatchLabels: t.getStoreLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: t.getRuleMeta(name),
					Spec: corev1.PodSpec{
						Volumes:        nil,
						InitContainers: nil,
						Containers: []corev1.Container{
							{
								Name:                     "",
								Image:                    "",
								Command:                  nil,
								Args:                     nil,
								WorkingDir:               "",
								Ports:                    nil,
								EnvFrom:                  nil,
								Env:                      nil,
								Resources:                corev1.ResourceRequirements{},
								VolumeMounts:             nil,
								VolumeDevices:            nil,
								LivenessProbe:            nil,
								ReadinessProbe:           nil,
								StartupProbe:             nil,
								Lifecycle:                nil,
								TerminationMessagePath:   "",
								TerminationMessagePolicy: "",
								ImagePullPolicy:          "",
								SecurityContext:          nil,
								Stdin:                    false,
								StdinOnce:                false,
								TTY:                      false,
							},
						},
						EphemeralContainers:           nil,
						RestartPolicy:                 "",
						TerminationGracePeriodSeconds: nil,
						ActiveDeadlineSeconds:         nil,
						DNSPolicy:                     "",
						NodeSelector:                  nil,
						ServiceAccountName:            "",
						DeprecatedServiceAccount:      "",
						AutomountServiceAccountToken:  nil,
						NodeName:                      "",
						HostNetwork:                   false,
						HostPID:                       false,
						HostIPC:                       false,
						ShareProcessNamespace:         nil,
						SecurityContext:               nil,
						ImagePullSecrets:              nil,
						Hostname:                      "",
						Subdomain:                     "",
						Affinity:                      nil,
						SchedulerName:                 "",
						Tolerations:                   nil,
						HostAliases:                   nil,
						PriorityClassName:             "",
						Priority:                      nil,
						DNSConfig:                     nil,
						ReadinessGates:                nil,
						RuntimeClassName:              nil,
						EnableServiceLinks:            nil,
						PreemptionPolicy:              nil,
						Overhead:                      nil,
						TopologySpreadConstraints:     nil,
					},
				},
			},
		}
		return statefulset, reconciler.StatePresent, nil
	}
	delete := &appsv1.StatefulSet{
		ObjectMeta: t.getObjectMeta(name),
	}
	return delete, reconciler.StateAbsent, nil
}
