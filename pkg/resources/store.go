package resources

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (t *ThanosComponentReconciler) storeService() (runtime.Object, reconciler.DesiredState, error) {
	name := "store-service"
	namespace := t.Thanos.Namespace
	if t.Thanos.Spec.StoreGateway != nil {
		store := t.Thanos.Spec.StoreGateway.DeepCopy()
		storeService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
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
				Selector: map[string]string{
					nameLabel:      "store",
					managedByLabel: t.Thanos.Name,
				},
				ClusterIP: corev1.ClusterIPNone,
				Type:      corev1.ServiceTypeClusterIP,
			},
		}
		return storeService, reconciler.StatePresent, nil

	}
	delete := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return delete, reconciler.StateAbsent, nil
}

func (t *ThanosComponentReconciler) storeDeployment() (runtime.Object, reconciler.DesiredState, error) {
	name := "store-deployment"
	namespace := t.Thanos.Namespace
	var objectStore *v1alpha1.ObjectStore
	if t.Thanos.Spec.ObjectStore != nil {
		if *t.Thanos.Spec.ObjectStore != "" {
			for _, o := range t.ObjectSores {
				if o.Name == *t.Thanos.Spec.ObjectStore {
					objectStore = o.DeepCopy()
				}
			}
			if objectStore == nil {
				return nil, nil, fmt.Errorf("unknown ObjectStore reference %q", *t.Thanos.Spec.ObjectStore)
			}
		} else if len(t.ObjectSores) > 0 {
			//TODO choose default better
			objectStore = t.ObjectSores[0].DeepCopy()
		} else {
			return nil, nil, fmt.Errorf("missing ObjectStore configuration for reference %q", *t.Thanos.Spec.ObjectStore)
		}
	} else {
		return nil, nil, fmt.Errorf("missing ObjectStore reference for Store")
	}

	if t.Thanos.Spec.StoreGateway != nil {
		store := t.Thanos.Spec.StoreGateway.DeepCopy()
		err := mergo.Merge(store, v1alpha1.DefaultStoreGateway) //TODO default Store
		if err != nil {
			return nil, nil, err
		}
		var deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace,
				Labels:      store.Labels,
				Annotations: store.Annotations,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						nameLabel:      "store",
						managedByLabel: t.Thanos.Name,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: "store",
						Labels: map[string]string{
							nameLabel:      "store",
							managedByLabel: t.Thanos.Name,
						},
						Annotations: store.Annotations,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "store",
								Image: fmt.Sprintf("%s:%s", store.Image.Repository, store.Image.Tag),
								Args: []string{
									"store",
									fmt.Sprintf("--grpc-address=%s", store.GRPCAddress),
									fmt.Sprintf("--http-address=%s", store.HttpAddress),
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", objectStore.Spec.Config.MountFrom.SecretKeyRef.Key),
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: GetPort(store.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: GetPort(store.GRPCAddress),
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
								ImagePullPolicy: corev1.PullIfNotPresent,
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "objectstore-secret",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: objectStore.Spec.Config.MountFrom.SecretKeyRef.Name,
									},
								},
							},
						},
					},
				},
			},
		}
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
