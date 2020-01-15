package resources

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// RunReconcilers run component reconcilers one by one and stop on first error.
func RunReconcilers(reconcilers []ComponentReconciler) (ctrl.Result, error) {
	for _, reconciler := range reconcilers {
		result, err := reconciler()
		if err != nil {
			return reconcile.Result{}, err
		}
		if result != nil {
			// short circuit if requested explicitly
			return *result, nil
		}
	}

	return ctrl.Result{}, nil
}
