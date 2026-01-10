---
linkTitle: "tcloud kubernetes nodepools delete"
title: "kubernetes nodepools delete"
slug: tcloud_kubernetes_nodepools_delete
url: /docs/tcloud/kubernetes/nodepools_delete/
weight: 9965
cascade:
  type: docs
---
## tcloud kubernetes nodepools delete

Delete a Kubernetes node pool

### Synopsis

Delete a node pool from a Kubernetes cluster.

This command will delete the node pool and all nodes associated with it.
The cluster must be in a ready state to delete node pools.

Examples:
  # Delete a node pool
  tcloud kubernetes nodepools delete --cluster my-cluster --name worker-pool

  # Delete a node pool and wait for completion
  tcloud kubernetes nodepools delete --cluster my-cluster --name worker-pool --wait

  # Delete a node pool without confirmation
  tcloud kubernetes nodepools delete --cluster my-cluster --name worker-pool --force

```
tcloud kubernetes nodepools delete [flags]
```

### Options

```
      --cluster string   Cluster identity, name, or slug (required)
      --force            Skip confirmation prompt
  -h, --help             help for delete
      --name string      Node pool name, identity, or slug (required)
      --wait             Wait for the node pool to be deleted before returning
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

* [tcloud kubernetes nodepools](/docs/tcloud/kubernetes/nodepools/)	 - Manage Kubernetes NodePools

