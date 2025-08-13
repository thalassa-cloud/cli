---
date: 2025-08-14T00:09:06+02:00
linkTitle: "tcloud kubernetes nodepools"
title: "kubernetes nodepools"
slug: tcloud_kubernetes_nodepools
url: /docs/tcloud/tcloud_kubernetes_nodepools/
weight: 9971
---
## tcloud kubernetes nodepools

Manage Kubernetes NodePools

### Examples

```
  # List all nodepools in a cluster
  tcloud kubernetes nodepools list my-cluster

  # Create a new nodepool
  tcloud kubernetes nodepools create my-cluster --name worker-pool --size 3

  # Delete a nodepool
  tcloud kubernetes nodepools delete my-cluster worker-pool
```

### Options

```
  -h, --help   help for nodepools
```

### Options inherited from parent commands

```
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
* [tcloud kubernetes nodepools list](/docs/tcloud/tcloud_kubernetes_nodepools_list/)	 - Kubernetes Cluster NodePool list

