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
)

func GetPort(address string) int32 {
	port, err := strconv.Atoi(strings.Split(address, ":")[1])
	if err != nil {
		return 0
	}
	return int32(port)
}

func (t *ThanosComponentReconciler) getStoreEndpoints() []string {
	var endpoints []string
	if t.Thanos.Spec.StoreGateway != nil {
		endpoints = append(endpoints, "store-service")
	}
	return endpoints
}

func (t *ThanosComponentReconciler) setQueryArgs(args []string) []string {
	query := t.Thanos.Spec.Query.DeepCopy()
	x := getArgs(query)
	fmt.Println(x)
	if query.LogLevel != "" {
		args = append(args, fmt.Sprintf("--log.level=%s", query.LogLevel))
	}
	if query.LogFormat != "" {
		args = append(args, fmt.Sprintf("--log.format=%s", query.LogFormat))
	}
	if query.GRPCGracePeriod != "" {
		args = append(args, fmt.Sprintf("--grpc-grace-period=%s", query.GRPCGracePeriod))
	}
	if query.WebRoutePrefix != "" {
		args = append(args, fmt.Sprintf("--web.route-prefix=%s", query.WebRoutePrefix))
	}
	if query.WebExternalPrefix != "" {
		args = append(args, fmt.Sprintf("--web.external-prefix=%s", query.WebExternalPrefix))
	}
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
	name := "query-deployment"
	namespace := t.Thanos.Namespace
	if t.Thanos.Spec.Query != nil {
		query := t.Thanos.Spec.Query.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace,
				Labels:      query.Labels,
				Annotations: query.Annotations,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "query"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "query",
						Labels:      map[string]string{"app": "query"},
						Annotations: query.Annotations,
					},
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
								ImagePullPolicy: corev1.PullIfNotPresent,
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return delete, reconciler.StateAbsent, nil
}
