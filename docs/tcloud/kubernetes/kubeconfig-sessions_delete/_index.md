---
linkTitle: "tcloud kubernetes kubeconfig-sessions delete"
title: "kubernetes kubeconfig-sessions delete"
slug: tcloud_kubernetes_kubeconfig-sessions_delete
url: /docs/tcloud/kubernetes/kubeconfig-sessions_delete/
weight: 9831
cascade:
  type: docs
---
## tcloud kubernetes kubeconfig-sessions delete

Delete a kubeconfig session

### Synopsis

Revoke a kubeconfig session for a Kubernetes cluster.

Provide the cluster as the first argument or with --cluster, and the session identity as the final argument.

```
tcloud kubernetes kubeconfig-sessions delete [cluster] <session> [flags]
```

### Options

```
      --cluster string   Cluster identity, name, or slug
      --force            Skip the confirmation prompt and delete
  -h, --help             help for delete
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

* [tcloud kubernetes kubeconfig-sessions](/docs/tcloud/kubernetes/kubeconfig-sessions/)	 - Manage Kubernetes kubeconfig sessions

