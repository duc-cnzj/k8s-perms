# k8s perms

> get all permission verbs in k8s.

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

```bash
# output
namespace 'default', kubeconfig '/Users/duc/goMod/checkmyk8sperm/rancherconfig'

core
- pods
    - delete
      authorization.k8s.io
- selfsubjectaccessreviews
    - create
- selfsubjectrulesreviews
    - create
```
