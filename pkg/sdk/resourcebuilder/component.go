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
	"fmt"
	"io/ioutil"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/static/gen/crds"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/static/gen/rbac"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

// +kubebuilder:object:generate=true

type ComponentConfig struct {
	Disabled              bool                 `json:"disabled,omitempty"`
	Name                  string               `json:"name,omitempty"`
	Namespace             string               `json:"namespace,omitempty"`
	MetaOverrides         *types.MetaBase      `json:"metaOverrides,omitempty"`
	WorkloadMetaOverrides *types.MetaBase      `json:"workloadMetaOverrides,omitempty"`
	WorkloadOverrides     *types.PodSpecBase   `json:"workloadOverrides,omitempty"`
	ContainerOverrides    *types.ContainerBase `json:"containerOverrides,omitempty"`
}

func(c *ComponentConfig) build(fn func(bool, *ComponentConfig) (runtime.Object, reconciler.DesiredState, error)) reconciler.ResourceBuilder {
	return reconciler.ResourceBuilder(func() (runtime.Object, reconciler.DesiredState, error) {
		return fn(c == nil || c.Disabled, c)
	})
}

func ResourceBuilders(object interface{}) []reconciler.ResourceBuilder {
	config := &ComponentConfig{
		Disabled: true,
	}
	if object != nil {
		config = object.(*ComponentConfig)
	}
	if config.Name == "" {
		config.Name = "thanos-operator"
	}
	if config.Namespace == "" {
		config.Namespace = "default"
	}
	resources := []reconciler.ResourceBuilder{
		config.build(Operator),
		config.build(ClusterRole),
		config.build(ClusterRoleBinding),
		config.build(ServiceAccount),
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(config, v1alpha1.GroupVersion.Group, "objectstores")
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(config, v1alpha1.GroupVersion.Group, "thanos")
		},
		func() (runtime.Object, reconciler.DesiredState, error) {
			return CRD(config, v1alpha1.GroupVersion.Group, "storeendpoints")
		},
	}
	return resources
}

func SetupWithBuilder(builder *builder.Builder) {
	builder.Owns(&appsv1.Deployment{})
}

func CRD(config *ComponentConfig, group string, kind string) (runtime.Object, reconciler.DesiredState, error) {
	crd := &v1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", kind, group),
		},
	}
	if config.Disabled {
		return crd, reconciler.StateAbsent, nil
	}
	crdFile, err := crds.Root.Open(fmt.Sprintf("/%s_%s.yaml", group, kind))
	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to open %s crd", kind)
	}
	bytes, err := ioutil.ReadAll(crdFile)
	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to read %s crd", kind)
	}

	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)

	_, _, err = serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, scheme, scheme, serializer.SerializerOptions{
		Yaml: true,
	}).Decode(bytes, &schema.GroupVersionKind{}, crd)

	if err != nil {
		return nil, nil, errors.WrapIff(err, "failed to unmarshal %s crd", kind)
	}

	// clear the TypeMeta to avoid objectmatcher diffing on it every time,
	// because the current object coming from the API Server will not have TypeMeta set
	crd.TypeMeta.Kind = ""
	crd.TypeMeta.APIVersion = ""

	return crd, reconciler.DesiredStateHook(func(object runtime.Object) error {
		current := object.(*v1beta1.CustomResourceDefinition)
		// simply copy the existing status over, so that we don't diff because of it
		crd.Status = current.Status
		return nil
	}), nil
}

func Operator(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: config.MetaOverrides.Merge(config.objectMeta()),
	}
	if disabled {
		return deployment, reconciler.StateAbsent, nil
	}
	deployment.Spec = appsv1.DeploymentSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: config.WorkloadMetaOverrides.Merge(v1.ObjectMeta{
				Labels: config.labelSelector(),
			}),
			Spec: config.WorkloadOverrides.Override(corev1.PodSpec{
				ServiceAccountName: config.objectMeta().Name,
				Containers: []corev1.Container{
					config.ContainerOverrides.Override(corev1.Container{
						Name:    "thanos-operator",
						Image:   "banzaicloud/thanos-operator",
						Command: []string{"/manager"},
						//Args:    []string{"--enable-leader-election"},
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
			MatchLabels: config.labelSelector(),
		},
	}
	return deployment, reconciler.StatePresent, nil
}

func ServiceAccount(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: config.MetaOverrides.Merge(config.objectMeta()),
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

func ClusterRoleBinding(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	rb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: config.MetaOverrides.Merge(config.clusterObjectMeta()),
	}

	if disabled {
		return rb, reconciler.StateAbsent, nil
	}

	sa := config.objectMeta().Name
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		sa = config.WorkloadOverrides.ServiceAccountName
	}

	rb.Subjects = []rbacv1.Subject{
		{
			Kind:      rbacv1.ServiceAccountKind,
			Name:      sa,
			Namespace: config.Namespace,
		},
	}
	rb.RoleRef = rbacv1.RoleRef{
		APIGroup: rbacv1.GroupName,
		Kind:     "ClusterRole",
		Name:     config.objectMeta().Name,
	}

	return rb, reconciler.StatePresent, nil
}

func ClusterRole(disabled bool, config *ComponentConfig) (runtime.Object, reconciler.DesiredState, error) {
	role := &rbacv1.ClusterRole{
		ObjectMeta: config.MetaOverrides.Merge(config.clusterObjectMeta()),
	}
	if disabled {
		return role, reconciler.StateAbsent, nil
	}
	// remove internal sa in case an externally provided service account is used
	if config.WorkloadOverrides != nil && config.WorkloadOverrides.ServiceAccountName != "" {
		return role, reconciler.StateAbsent, nil
	}
	roleFile, err := rbac.Root.Open("/role.yaml")
	if err != nil {
		return nil, nil, errors.WrapIf(err, "failed to open role.yaml")
	}
	roleAsByte, err := ioutil.ReadAll(roleFile)
	if err != nil {
		return nil, nil, err
	}
	scheme := runtime.NewScheme()
	rbacv1.AddToScheme(scheme)
	_, _, err = serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, scheme, scheme, serializer.SerializerOptions{
		Yaml: true,
	}).Decode(roleAsByte, &schema.GroupVersionKind{}, role)

	// overwrite the objectmeta that has been read from file
	role.ObjectMeta = config.MetaOverrides.Merge(config.clusterObjectMeta())
	role.TypeMeta.Kind = ""
	role.TypeMeta.APIVersion = ""

	return role, reconciler.StatePresent, err
}

func (c *ComponentConfig) objectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:      c.Name,
		Namespace: c.Namespace,
		Labels:    c.labelSelector(),
	}
	return meta
}

func (c *ComponentConfig) clusterObjectMeta() v1.ObjectMeta {
	meta := v1.ObjectMeta{
		Name:   c.Name,
		Labels: c.labelSelector(),
	}
	return meta
}

func (c *ComponentConfig) labelSelector() map[string]string {
	return map[string]string{
		"banzaicloud.io/operator": c.Name,
	}
}
