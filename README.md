# ca-controller-for-strimzi

> :warning: **This project is not in cooperation with Strimzi**: I asked one of the Strimzi core maintainer about
> contributing this controller upstream but it makes more sense to implement the logic in Strimzi itself.

When dealing with [strimzi-kafka-operator](https://github.com/strimzi/strimzi-kafka-operator/), it is possible to use
your own CA for cluster and clients.
Strimzi requires the CA key split from the rest of a Kubernetes secret of type TLS.
To avoid doing this manually e.g. when using [cert-manager](https://cert-manager.io/) for managing the CAs, this
controller can be used...

Please check the `examples/example-ca.yaml` how to use the controller after deploying it and using it with cert-manager but
it also works with normal Kubernetes secrets of type TLS.

## Install the controller via Helm

```shell
helm repo add sebastiangaiser-ca-controller-for-strimzi https://sebastiangaiser.github.io/ca-controller-for-strimzi
helm repo update
helm upgrade --install -n kafka ca-controller-for-strimzi sebastiangaiser-ca-controller-for-strimzi/ca-controller-for-strimzi
```

## See the controller in action (using local kind cluster)

```shell
# create cluster
kind create cluster
# install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml
# bootstrap strimzi in kafka namespace
kubectl create namespace kafka
helm install -n kafka strimzi-cluster-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator
# create certificates
kubectl apply -f examples/example-ca.yaml
# install controller
helm repo add sebastiangaiser-ca-controller-for-strimzi https://sebastiangaiser.github.io/ca-controller-for-strimzi
helm repo update
helm upgrade --install -n kafka ca-controller-for-strimzi sebastiangaiser-ca-controller-for-strimzi/ca-controller-for-strimzi
# create Kafka cluster
kubectl apply -f examples/kafka.yaml -n kafka
```
