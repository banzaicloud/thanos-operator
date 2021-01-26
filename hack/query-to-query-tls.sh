#kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml

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
  name: peer-ca
spec:
  commonName: peer-ca
  secretName: peer-ca
  isCA: true
  issuerRef:
    name: selfsigned
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: observer-ca
spec:
  commonName: observer-ca
  secretName: observer-ca
  isCA: true
  issuerRef:
    name: selfsigned
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: peer
spec:
  ca:
    secretName: peer-ca
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: observer
spec:
  ca:
    secretName: observer-ca
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: peer-tls
spec:
  secretName: peer-tls
  dnsNames:
  - peer-query.default.svc.cluster.local
  - 127.0.0.1.xip.io
  issuerRef:
    name: peer
  usages:
  - server auth
EOF

cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: observer-tls
spec:
  commonName: observer
  secretName: observer-tls
  issuerRef:
    name: observer
  usages:
  - client auth
EOF

mkdir -p certs/observer certs/peer

kubectl view-secret peer-tls ca.crt > certs/observer/ca.crt
kubectl view-secret observer-tls tls.crt > certs/observer/tls.crt
kubectl view-secret observer-tls tls.key > certs/observer/tls.key

kubectl view-secret observer-tls ca.crt > certs/peer/ca.crt
kubectl view-secret peer-tls tls.crt > certs/peer/tls.crt
kubectl view-secret peer-tls tls.key > certs/peer/tls.key

kubectl create secret generic peer-tls-with-observer-ca --from-file certs/peer
kubectl create secret generic observer-tls-with-peer-ca --from-file certs/observer

# Aggregator query
cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: observer
spec:
  # Add fields here
  query:
    GRPCClientCertificate: "observer-tls-with-peer-ca"
  queryDiscovery: true
EOF

# Peer query
cat <<EOF | kubectl replace -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: peer
spec:
  # Add fields here
  query:
    GRPCServerCertificate: "peer-tls-with-observer-ca"
EOF