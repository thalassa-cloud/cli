---
linkTitle: "tcloud kubernetes connect"
title: "kubernetes connect"
slug: tcloud_kubernetes_connect
url: /docs/tcloud/kubernetes/connect/
weight: 9974
cascade:
  type: docs
---
## tcloud kubernetes connect

Connect your shell to the Kubernetes Cluster

```
tcloud kubernetes connect [flags]
```

### Options

```
  -h, --help                     help for connect
      --kubeconfig-path string   path to the kubeconfig file
      --temp                     use a temporary kubeconfig file (default true)
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes](/docs/tcloud/tcloud_kubernetes/)	 - Manage Kubernetes clusters, node pools and more services related to Kubernetes

