# namespace-populator

A Kubernetes controller to auto generate ressources based on namespace creation

## Goal

When a user create a new namespace, some default resources may be spawned in ou out this namespace.

## How

This controller listen on namespace creation and is configured with configmaps.

Configmaps are 1 ressource (yaml) by files. Resources can be changed with go template on namespace creation.

Some options can also be add as annotation to filter on namespaces.

## Use cases

* default LimitRange object
* default network policies
* monitoring by namespace
* ...

## [Examples](/examples)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  labels:
    app: namespace-populator
  annotations:
    # this annotation can be used to select which namespace to trigger
    namespace-populator.barpilot.io/selector: "create=nginx"
data:
  pod.yaml: |
    apiVersion: v1
    kind: Pod
    metadata:
      name: nginx
      namespace: {{ .Name }}
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

## Limitations

* Resources aren't removed on namespace deletion (this can be added later #1 )
* Resources are created but not updated
* Resources deleted are recreated (even if no new namespaces are created, controller process namespaces each 30 sec)
* Namespace selection can only be done by labels. This can be bypass using `if` on `.Name` in go template.
