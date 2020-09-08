package query_frontend

import (
	"fmt"
	"sort"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(reconciler *resources.ThanosComponentReconciler) *QueryFrontend {
	return &QueryFrontend{
		ThanosComponentReconciler: reconciler,
	}
}

type QueryFrontend struct {
	*resources.ThanosComponentReconciler
}

func (q *QueryFrontend) Reconcile() (*reconcile.Result, error) {
	if q.Thanos.Spec.QueryFrontend != nil {
		err := mergo.Merge(q.Thanos.Spec.QueryFrontend, v1alpha1.DefaultQueryFrontend)
		if err != nil {
			return nil, err
		}
	}
	return q.ReconcileResources(
		[]resources.Resource{
			q.deployment,
			q.service,
			q.serviceMonitor,
			q.ingressHTTP,
			q.ingressGRPC,
		})
}

func (q *QueryFrontend) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: v1alpha1.QueryFrontendName,
	}.Merge(
		q.GetCommonLabels(),
	)
	return labels
}

func (q *QueryFrontend) getName(suffix ...string) string {
	name := v1alpha1.QueryFrontendName
	if len(suffix) > 0 {
		name = name + "-" + suffix[0]
	}
	return q.QualifiedName(name)
}

func (q *QueryFrontend) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc.cluster.local", q.getName(), q.Thanos.Namespace)
}

func (q *QueryFrontend) getMeta(name string, params ...string) metav1.ObjectMeta {
	namespace := ""
	if len(params) > 0 {
		namespace = params[0]
	}
	meta := q.GetObjectMeta(name, namespace)
	meta.Labels = q.getLabels()
	return meta
}

func (q *QueryFrontend) setArgs(originArgs []string) []string {
	query := q.Thanos.Spec.QueryFrontend.DeepCopy()
	// Get args from the tags
	args := resources.GetArgs(query)

	if query.GRPCClientCertificate != "" {
		args = append(args, "--grpc-client-tls-secure")
		args = append(args, fmt.Sprintf("--grpc-client-tls-cert=%s/%s", clientCertMountPath, "tls.crt"))
		args = append(args, fmt.Sprintf("--grpc-client-tls-key=%s/%s", clientCertMountPath, "tls.key"))
		args = append(args, fmt.Sprintf("--grpc-client-tls-ca=%s/%s", clientCertMountPath, "ca.crt"))
		args = append(args, "--grpc-client-server-name=example.com") //TODO this is dummy now
	}
	if query.GRPCServerCertificate != "" {
		args = append(args, fmt.Sprintf("--grpc-server-tls-cert=%s/%s", serverCertMountPath, "tls.crt"))
		args = append(args, fmt.Sprintf("--grpc-server-tls-key=%s/%s", serverCertMountPath, "tls.key"))
		args = append(args, fmt.Sprintf("--grpc-server-tls-client-ca=%s/%s", serverCertMountPath, "ca.crt"))
	}
	// Handle special args
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
	args = append(args, q.getStoreEndpoints()...)

	// Sort generated args to prevent accidental diffs
	sort.Strings(args)
	// Concat original and computed args
	finalArgs := append(originArgs, args...)
	return finalArgs
}
