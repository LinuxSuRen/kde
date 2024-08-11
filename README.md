# kde
Kubernetes based remote IDE

## Setup Project

I setup this project using the following commands:

* kubebuilder init --plugins=go.kubebuilder.io/v4
* kubebuilder create api --group linuxsuren.github.io --version v1alpha1 --kind DevSpace

## Install OpenEBS

```shell
kubectl apply -f https://openebs.github.io/charts/openebs-operator.yaml
```

## Install Nginx Ingress

```shell
helm upgrade --install ingress-nginx ingress-nginx \
  --set controller.image.registry=k8s.m.daocloud.io \
  --set controller.admissionWebhooks.patch.image.registry=k8s.m.daocloud.io \
  --set controller.opentelemetry.image.registry=k8s.m.daocloud.io \
  --set defaultBackend.image.registry=k8s.m.daocloud.io \
  --set defaultBackend.service.type=NodePort \
  --set controller.service.type=NodePort \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```
