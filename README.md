# k8s perms

> get all permission verbs in k8s.

```bash
# k8s-perms -n default -kubeconfig ~/.kube/config

# output
namespace 'default', kubeconfig '~/.kube/config'

core
  - pods
    - delete
authorization.k8s.io
  - selfsubjectaccessreviews
    - create
  - selfsubjectrulesreviews
    - create
```


## Installation

```bash
go install github.com/duc-cnzj/k8s-perms@latest
```

## Usage

```bash
k8s-perms -n default -kubeconfig ~/.kube/config

or 

export KUBECONFIG=~/.kube/config

k8s-perms -n default
```