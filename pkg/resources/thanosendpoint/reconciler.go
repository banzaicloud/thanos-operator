package thanosendpoint

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log logr.Logger
	resourceReconciler reconciler.ResourceReconciler
}

func NewReconciler(logger logr.Logger, reconciler reconciler.ResourceReconciler) *Reconciler {
	return &Reconciler{
		log: logger,
		resourceReconciler: reconciler,
	}
}

func (r Reconciler) Reconcile() (*reconcile.Result, error) {
	r.log.Info("reconciling all day")
	return nil, nil
}

