# ca-controller-for-strimzi

> :warning: **This project is not in cooperation with Strimzi**: I asked one of the Strimzi core maintainer about
> contributing this controller upstream but it makes more sense to implement the logic in Strimzi itself.

When dealing with [strimzi-kafka-operator](https://github.com/strimzi/strimzi-kafka-operator/), it is possible to use
your own CA for cluster and clients.
Strimzi requires the CA key split from the rest of a Kubernetes secret of type TLS.
To avoid doing this manually e.g. when using [cert-manager](https://cert-manager.io/) for managing the CAs, this
controller can be used...

Please check the `example-ca.yaml` how to use the controller after deploying it and using it with cert-manager but
it also works with normal Kubernetes secrets of type TLS.
