# tkn

This task performs operations on Tekton resources using
[`tkn`](https://github.com/tektoncd/cli).

## Install the Task

```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/tkn/0.1/tkn.yaml
```

## Parameters

name      | description                                 | default
--------- | ------------------------------------------- | -------
tkn-image | `tkn` CLI container image to run this task. | gcr.io/tekton-releases/dogfooding/tkn
ARGS      | The arguments to pass to the `tkn` CLI.     | No default. Required.

### Platforms

The Task can be run on `linux/amd64`, `linux/s390x` and `linux/ppc64le` platforms.
