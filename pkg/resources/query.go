package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const queryName = "query"

func GetPort(address string) int32 {
	res := strings.Split(address, ":")
	if len(res) > 1 {
		port, err := strconv.Atoi(res[1])
		if err != nil {
			return 0
		}
		return int32(port)
	}
	return 0
}

func (t *ThanosComponentReconciler) getQueryLabels() Labels {
	return Labels{
		nameLabel: queryName,
	}.merge(
		t.Thanos.Spec.Query.Labels,
		t.getCommonLabels(),
	)
}

func (t *ThanosComponentReconciler) getQueryMeta(name string) metav1.ObjectMeta {
	meta := t.getObjectMeta(name)
	meta.Labels = t.getQueryLabels()
	meta.Annotations = t.Thanos.Spec.Query.Annotations
	return meta
}

func (t *ThanosComponentReconciler) getStoreEndpoints() []string {
	var endpoints []string
	if t.Thanos.Spec.StoreGateway != nil {
		endpoints = append(endpoints, t.qualifiedName(storeName))
	}
	return endpoints
}

func (t *ThanosComponentReconciler) setQueryArgs(args []string) []string {
	query := t.Thanos.Spec.Query.DeepCopy()
	args = append(args, getArgs(query)...)
	if query.QueryReplicaLabels != nil {
		for _, l := range query.QueryReplicaLabels {
			args = append(args, fmt.Sprintf("--query.replica-label=%s", l))
		}
	}
	if query.SelectorLabels != nil {
		for k, v := range query.SelectorLabels {
			args = append(args, fmt.Sprintf("--selector-label=%s=%s", k, v))
		}
	}
	// Add discovery args
	if t.Thanos.Spec.ThanosDiscovery != nil {
		for _, s := range t.getStoreEndpoints() {
			args = append(args, fmt.Sprintf("--store=dnssrvnoa+_grpc._tcp.%s.%s.svc.cluster.local", s, t.Thanos.Namespace))
		}
	}
	return args
}

func (t *ThanosComponentReconciler) queryDeployment() (runtime.Object, reconciler.DesiredState, error) {
	if t.Thanos.Spec.Query != nil {
		query := t.Thanos.Spec.Query.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: t.getQueryMeta(queryName),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: t.getQueryLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: t.getQueryMeta(queryName),
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
										ContainerPort: GetPort(query.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: GetPort(query.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								Resources:       query.Resources,
								ImagePullPolicy: query.Image.PullPolicy,
								LivenessProbe:   t.getCheck(GetPort(query.HttpAddress), healthCheckPath),
								ReadinessProbe:  t.getCheck(GetPort(query.HttpAddress), readyCheckPath),
							},
						},
					},
				},
			},
		}
		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = t.setQueryArgs(deployment.Spec.Template.Spec.Containers[0].Args)
		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: t.getObjectMeta(queryName),
	}
	return delete, reconciler.StateAbsent, nil
}

func (t *ThanosComponentReconciler) queryService() (runtime.Object, reconciler.DesiredState, error) {
	if t.Thanos.Spec.StoreGateway != nil {
		store := t.Thanos.Spec.StoreGateway.DeepCopy()
		storeService := &corev1.Service{
			ObjectMeta: t.getQueryMeta(queryName),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "grpc",
						Protocol: corev1.ProtocolTCP,
						Port:     GetPort(store.GRPCAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "grpc",
						},
					},
					{
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						Port:     GetPort(store.HttpAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "http",
						},
					},
				},
				Selector: t.getQueryLabels(),
				Type:     corev1.ServiceTypeClusterIP,
			},
		}
		return storeService, reconciler.StatePresent, nil

	}
	delete := &corev1.Service{
		ObjectMeta: t.getQueryMeta(queryName),
	}
	return delete, reconciler.StateAbsent, nil
}
