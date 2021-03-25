#!/bin/bash
# start a cluster on aws

# install cert-manager

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml

# nginx ingress

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.44.0/deploy/static/provider/cloud/deploy.yaml

# install minio and prometheus

one-eye prometheus install --update --accept-license

# start the thanos operator on demand in a separate shell

make run &
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

# we need an alias to the endpoint to generate the cert with
export PEER_ENDPOINT=peer

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
kind: ThanosEndpoint
metadata:
  name: $PEER_ENDPOINT
spec:
  certificate: peer-tls
  caBundle: peer-tls
  queryOverrides:
    GRPCIngress:
      ingressOverrides:
        metadata:
          annotations:
            kubernetes.io/ingress.class: "nginx"
EOF

# wait for the ingress endpoint and register it as a cname in an externalname service

while true; do
  export INGRESS_ENDPOINT=$(kubectl get thanosendpoint ${PEER_ENDPOINT} -o jsonpath='{.status.endpointAddress}')
  [[ -n $INGRESS_ENDPOINT ]] && break
  echo -n "." && sleep 1
done

arrIN=(${INGRESS_ENDPOINT//:/ })

kubectl create service externalname $PEER_ENDPOINT --external-name ${arrIN[0]} || true

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
