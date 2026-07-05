---
linkTitle: "tcloud kubernetes machines list"
title: "kubernetes machines list"
slug: tcloud_kubernetes_machines_list
url: /docs/tcloud/kubernetes/machines_list/
weight: 9884
cascade:
  type: docs
---
## tcloud kubernetes machines list

List machines in a Kubernetes cluster

### Synopsis

Lists worker machines (nodes) across node pools in the given cluster.

```
tcloud kubernetes machines list [flags]
```

### Options

```
      --cluster string    Cluster identity, name, or slug
  -h, --help              help for list
      --no-header         Do not print the header
      --nodepool string   Filter by node pool identity, name, or slug
      --show-exact-time   Show exact time instead of relative time
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

* [tcloud kubernetes machines](/docs/tcloud/kubernetes/machines/)	 - List and manage Kubernetes cluster machines

