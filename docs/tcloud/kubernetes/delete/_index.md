---
linkTitle: "tcloud kubernetes delete"
title: "kubernetes delete"
slug: tcloud_kubernetes_delete
url: /docs/tcloud/kubernetes/delete/
weight: 9970
cascade:
  type: docs
---
## tcloud kubernetes delete

Delete a Kubernetes cluster

### Synopsis

Delete a Kubernetes cluster and all associated resources.

This command will delete the cluster and all node pools, nodes, and other resources
associated with it. This operation cannot be undone.

Examples:
  # Delete a cluster
  tcloud kubernetes delete my-cluster

  # Delete a cluster and wait for completion
  tcloud kubernetes delete my-cluster --wait

  # Delete a cluster without confirmation
  tcloud kubernetes delete my-cluster --force

```
tcloud kubernetes delete <cluster> [flags]
```

### Options

```
      --force   Skip confirmation prompt
  -h, --help    help for delete
      --wait    Wait for the cluster to be deleted before returning
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

