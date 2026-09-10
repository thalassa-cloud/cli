---
linkTitle: "tcloud kubernetes kubeconfig-sessions"
title: "kubernetes kubeconfig-sessions"
slug: tcloud_kubernetes_kubeconfig-sessions
url: /docs/tcloud/kubernetes/kubeconfig-sessions/
weight: 9829
cascade:
  type: docs
---
## tcloud kubernetes kubeconfig-sessions

Manage Kubernetes kubeconfig sessions

### Synopsis

List and revoke active kubeconfig sessions for a Kubernetes cluster.

### Examples

```
  # List kubeconfig sessions for a cluster
  tcloud kubernetes kubeconfig-sessions list my-cluster

  # Delete a kubeconfig session
  tcloud kubernetes kubeconfig-sessions delete my-cluster sess-abc123 --force
```

### Options

```
  -h, --help   help for kubeconfig-sessions
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes](/docs/tcloud/tcloud_kubernetes/)	 - Manage Kubernetes clusters, node pools and more services related to Kubernetes
* [tcloud kubernetes kubeconfig-sessions delete](/docs/tcloud/kubernetes/kubeconfig-sessions_delete/)	 - Delete a kubeconfig session
* [tcloud kubernetes kubeconfig-sessions list](/docs/tcloud/kubernetes/kubeconfig-sessions_list/)	 - List kubeconfig sessions for a Kubernetes cluster

