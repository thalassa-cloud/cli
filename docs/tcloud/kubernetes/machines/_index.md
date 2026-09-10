---
linkTitle: "tcloud kubernetes machines"
title: "kubernetes machines"
slug: tcloud_kubernetes_machines
url: /docs/tcloud/kubernetes/machines/
weight: 9826
cascade:
  type: docs
---
## tcloud kubernetes machines

List and manage Kubernetes cluster machines

### Synopsis

Commands for machines (nodes) that belong to Kubernetes node pools within a cluster.

### Examples

```
  # List all machines in a cluster
  tcloud kubernetes machines list --cluster my-cluster

  # List machines in a specific node pool
  tcloud kubernetes machines list --cluster my-cluster --nodepool worker
```

### Options

```
  -h, --help   help for machines
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
* [tcloud kubernetes machines list](/docs/tcloud/kubernetes/machines_list/)	 - List machines in a Kubernetes cluster

