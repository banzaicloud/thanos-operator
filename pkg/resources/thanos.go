package resources

import "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"

type ThanosReconciler struct {
	Query *v1alpha1.Query
}
