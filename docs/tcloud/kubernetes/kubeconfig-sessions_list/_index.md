---
linkTitle: "tcloud kubernetes kubeconfig-sessions list"
title: "kubernetes kubeconfig-sessions list"
slug: tcloud_kubernetes_kubeconfig-sessions_list
url: /docs/tcloud/kubernetes/kubeconfig-sessions_list/
weight: 9830
cascade:
  type: docs
---
## tcloud kubernetes kubeconfig-sessions list

List kubeconfig sessions for a Kubernetes cluster

```
tcloud kubernetes kubeconfig-sessions list [cluster] [flags]
```

### Options

```
      --cluster string   Cluster identity, name, or slug
      --exact-time       Show full timestamps instead of relative time
  -h, --help             help for list
      --no-header        Do not print the header
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

