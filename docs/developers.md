<p align="center"><img src="./img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

# Developers documentation

## Rapid local flow

This is best used in rapid development workflows, you will need two terminals for this.

1. Start thanos-operator locally on your laptop, it will target the current Kubernetes context:

    ```bash
    make run
    ```

2. In another terminal create a sample Thanos installation with Minio, Prometheus and Thanos CRDs:

    ```bash
    make install-thanos
    ```

3. Check that Thanos runs properly:
    ```bash
    kubectl port-forward service/thanos-sample-query 10902 &
    open http://localhost:10902/stores
    ```

## Containerized flow

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
