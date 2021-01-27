#!/bin/bash
# start a cluster on aws

# install cert-manager

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml

# install minio and prometheus

make install-minio install-prometheus

# start the thanos operator on demand in a separate shell

make run

# we need an alias to the endpoint to generate the cert with
export PEER_ENDPOINT=peer-endpoint

# create a selfsigned CA and a cert for communication
cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned
spec:
  selfSigned: {}
EOF

# this will be the single cert for both sides
cat <<EOF | kubectl apply -f-
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: peer-tls
spec:
  secretName: peer-tls
  commonName: peer-endpoint.cluster.notld
  dnsNames:
  - $PEER_ENDPOINT
  issuerRef:
    name: selfsigned
  usages:
  - server auth
  - client auth
EOF

# create a peer that will be routed out to the public internet through a grpc ingress

cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: peer
spec:
  query:
    GRPCIngress:
      ingressOverrides:
        metadata:
          annotations:
            nginx.ingress.kubernetes.io/backend-protocol: GRPC
            kubernetes.io/ingress.class: "nginx"
            nginx.ingress.kubernetes.io/auth-tls-verify-client: "on"
            nginx.ingress.kubernetes.io/auth-tls-secret: "default/peer-tls"
      certificate: peer-tls
      host: $PEER_ENDPOINT
EOF

cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: StoreEndpoint
metadata:
  name: peer
spec:
  thanos: peer
  selector: {}
EOF

# wait for the ingress endpoint and register it as a cname in an externalname service

while true; do
  export INGRESS_ENDPOINT=$(kubectl get ing peer-query-grpc -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
  [[ -n $INGRESS_ENDPOINT ]] && break
  echo -n "." && sleep 1
done

kubectl create service externalname $PEER_ENDPOINT --external-name $INGRESS_ENDPOINT

# create our central query instance that will connect to our external peer endpoint through the external-name service
cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: observer
spec:
  # Add fields here
  query:
    GRPCClientCertificate: "peer-tls"
EOF

cat <<EOF | kubectl apply -f-
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: StoreEndpoint
metadata:
  name: observer
spec:
  thanos: observer
  url: ${PEER_ENDPOINT}:443
EOF
