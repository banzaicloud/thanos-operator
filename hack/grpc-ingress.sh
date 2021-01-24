# Example: https://github.com/kubernetes/ingress-nginx/issues/5886#issuecomment-657923144

# Create a cluster using kind: https://kind.sigs.k8s.io/docs/user/ingress/#create-cluster
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

# Install ingress-nginx: https://kind.sigs.k8s.io/docs/user/ingress/#ingress-nginx
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

# install cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml

# install prometheus operator (CRDs only for this demo)
(cd /tmp && go get github.com/banzaicloud/k8s-yaml-filter)
helm template prometheus-community/kube-prometheus-stack --include-crds | k8s-yaml-filter -kind CustomResourceDefinition

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned
spec:
  selfSigned: {}
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: xip-tls
spec:
  secretName: xip-tls
  dnsNames:
  - 127.0.0.1.xip.io
  issuerRef:
    name: selfsigned
EOF

cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: query
spec:
  # Add fields here
  query:
    GRPCIngress:
      ingressOverrides:
        metadata:
          annotations:
            nginx.ingress.kubernetes.io/backend-protocol: GRPC
      certificate: xip-tls
      host: 127.0.0.1.xip.io
EOF