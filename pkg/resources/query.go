package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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

func (t *ThanosReconciler) queryDeployment() (runtime.Object, reconciler.DesiredState, error) {
	query := t.Query.DeepCopy()
	if t.Query != nil {
		err := mergo.Merge(query, v1alpha1.DefaultQuery)
		if err != nil {
			return nil, nil, err
		}
		query := t.Query
		var deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "",
				Namespace:   "",
				Labels:      query.Labels,
				Annotations: query.Annotations,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "query"},
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "",
						Namespace:       "",
						ResourceVersion: "",
						Labels:          map[string]string{"app": "query"},
						Annotations:     query.Annotations,
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "query",
								Image: fmt.Sprintf("%s:%s", query.Image.Repository, query.Image.Tag),
								Args: []string{
									"query",
									fmt.Sprintf("--grpc-address=%s", query.GRPCAddress),
									fmt.Sprintf("--http-address=%s", query.HttpAddress),
								},
								Ports: []v1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: GetPort(query.HttpAddress),
										Protocol:      "tcp",
									},
									{
										Name:          "grpc",
										ContainerPort: GetPort(query.GRPCAddress),
										Protocol:      "tcp",
									},
								},
								Resources:       query.Resources,
								ImagePullPolicy: v1.PullPolicy(query.Image.PullPolicy),
							},
						},
					},
				},
			},
		}
		return deployment, reconciler.StatePresent, nil
	}
	return nil, reconciler.StateAbsent, nil

}
