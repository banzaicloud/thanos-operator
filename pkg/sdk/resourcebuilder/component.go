// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resourcebuilder

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

const ThanosOperator = "thanos-operator"

// +kubebuilder:object:generate=true

type ComponentConfig struct {
	Disabled           bool                 `json:"disabled,omitempty"`
	WorkloadOverrides  *types.WorkloadBase  `json:"workloadOverrides,omitempty"`
	ContainerOverrides *types.ContainerBase `json:"containerOverrides,omitempty"`
}

func ResourceBuilders(object interface{}) []reconciler.ResourceBuilder {
	var config *ComponentConfig
	if object != nil {
		config = object.(*ComponentConfig)
	}
	disabled := (config == nil || config.Disabled)

	resources := []reconciler.ResourceBuilder{
		func() (runtime.Object, reconciler.DesiredState, error) {
			return Operator(disabled, config)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return Role(disabled, config)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return RoleBinding(disabled, config)
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return ServiceAccount(disabled, config)
		},
	}
	return resources
}

func SetupWithBuilder(builder *builder.Builder) {
	builder.Owns(&appsv1.Deployment{})
}

func Operator(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	if disabled {
		return &appsv1.Deployment{
			ObjectMeta: objectMeta(),
		}, reconciler.StateAbsent, nil
	}
	return &appsv1.Deployment{
		ObjectMeta: objectMeta(),
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: config.WorkloadOverrides.MergeMeta(v1.ObjectMeta{
					Labels: labelSelector(),
				}),
				Spec: config.WorkloadOverrides.OverrideSpec(corev1.PodSpec{
					ServiceAccountName: ThanosOperator,
					Containers: []corev1.Container{
						config.ContainerOverrides.Override(corev1.Container{
							Name:    "thanos-operator",
							Image:   "banzaicloud/thanos-operator",
							Command: []string{"/manager"},
							Args:    []string{"--enable-leader-election"},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("30Mi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("50m"),
									corev1.ResourceMemory: resource.MustParse("20Mi"),
								},
							},
						}),
					},
				}),
			},
			Selector: &v1.LabelSelector{
				MatchLabels: labelSelector(),
			},
		},
	}, reconciler.StatePresent, nil
}

func ServiceAccount(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:   ThanosOperator,
			Labels: labelSelector(),
		},
	}
	if disabled {
		return sa, reconciler.StateAbsent, nil
	}
	// remove internal sa in case an externally provided service account is used
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		return sa, reconciler.StateAbsent, nil
	}
	return sa, reconciler.StatePresent, nil
}

func RoleBinding(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: objectMeta(),
	}

	if disabled {
		return rb, reconciler.StateAbsent, nil
	}

	sa := ThanosOperator
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		sa = config.WorkloadOverrides.ServiceAccountName
	}

	rb.Subjects = []rbacv1.Subject{
		{
			Kind: rbacv1.ServiceAccountKind,
			Name: sa,
		},
	}
	rb.RoleRef = rbacv1.RoleRef{
		APIGroup: rbacv1.GroupName,
		Kind:     "Role",
		Name:     objectMeta().Name,
	}

	return rb, reconciler.StatePresent, nil
}

func Role(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	role := &rbacv1.Role{
		ObjectMeta: objectMeta(),
	}
	if disabled {
		return role, reconciler.StateAbsent, nil
	}
	// remove internal sa in case an externally provided service account is used
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		return role, reconciler.StateAbsent, nil
	}
	roleAsString := `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
rules:
- apiGroups:
  - ""
  resources:
  - services
  - configmaps
  - events
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
- apiGroups:
  - apps
  - extensions
  resources:
  - deployments
  - statefulsets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - objectstores
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - objectstores/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - storeendpoints
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - storeendpoints/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - thanos
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - monitoring.banzaicloud.io
  resources:
  - thanos/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - monitoring.coreos.com
  resources:
  - servicemonitors
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
`
	scheme := runtime.NewScheme()
	rbacv1.AddToScheme(scheme)
	_, _, err := serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, scheme, scheme, serializer.SerializerOptions{
		Yaml: true,
	}).Decode([]byte(roleAsString), &schema.GroupVersionKind{}, role)

	role.TypeMeta.Kind = ""
	role.TypeMeta.APIVersion = ""

	role.ObjectMeta = objectMeta()

	return role, reconciler.StatePresent, err
}


func objectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:   ThanosOperator,
		Labels: labelSelector(),
	}
	return meta
}

func labelSelector() map[string]string {
	return map[string]string{
		"banzaicloud.io/operator": ThanosOperator,
	}
}
