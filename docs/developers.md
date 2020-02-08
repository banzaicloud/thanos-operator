<p align="center"><img src="./img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

# Developers documentation

### Build & Push container
```
make docker-build IMG=banzaicloud/thanos-operator:latest
make docker-push IMG=banzaicloud/thanos-operator:latest
```

### Install operator the cluster
```
export KUBECONFIG="<kubeconfig location>"
make install
make deploy IMG=banzaicloud/thanos-operator:latest
```
