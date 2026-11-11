---
linkTitle: "tcloud kubernetes nodepools update"
title: "kubernetes nodepools update"
slug: tcloud_kubernetes_nodepools_update
url: /docs/tcloud/kubernetes/nodepools_update/
weight: 9805
cascade:
  type: docs
---
## tcloud kubernetes nodepools update

Update a Kubernetes node pool

### Synopsis

Update an existing node pool in a Kubernetes cluster.

Only fields specified with flags will be updated. All other fields will remain unchanged.

Examples:
  # Update node pool replicas
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5

  # Enable autoscaling
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --enable-autoscaling --min-nodes 2 --max-nodes 10

  # Update machine type
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --machine-type pgp-large

  # Update the Kubernetes version
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --kubernetes-version 1.32.2

  # Update multiple fields
  tcloud kubernetes nodepools update --cluster my-cluster --name worker --num-nodes 5 --enable-autohealing

To choose the next version within the Kubernetes version policy, use 'tcloud kubernetes nodepools upgrade'.

Note: Use 'tcloud kubernetes nodepools label' and 'tcloud kubernetes nodepools annotate' commands
to manage labels and annotations separately. Those values are resent with the update so they are not cleared.

```
tcloud kubernetes nodepools update [flags]
```

### Options

```
      --cluster string              Cluster identity, name, or slug (defaults to kubernetes.cluster in .thalassa)
      --enable-autohealing          Enable autohealing for the node pool
      --enable-autoscaling          Enable autoscaling for the node pool
  -h, --help                        help for update
      --kubernetes-version string   Kubernetes version for the node pool (identity, slug, name, or version)
      --machine-type string         Machine type for the node pool
      --max-nodes int               Maximum number of nodes (when autoscaling is enabled)
      --min-nodes int               Minimum number of nodes (when autoscaling is enabled)
      --name string                 Node pool name, identity, or slug (required)
      --node-taints strings         Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule'). Replaces existing taints.
      --num-nodes int               Number of nodes in the node pool (only when autoscaling is disabled)
      --security-groups strings     Security group identities to attach to node pool machines
      --upgrade-strategy string     Upgrade strategy: manual, auto, always, on-delete, inplace, or never
      --wait                        Wait for the node pool update to complete
      --wait-timeout duration       Maximum time to wait for the node pool update to complete (default 20m0s)
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes nodepools](/docs/tcloud/kubernetes/nodepools/)	 - Manage Kubernetes NodePools

