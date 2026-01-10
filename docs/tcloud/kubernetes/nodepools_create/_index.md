---
linkTitle: "tcloud kubernetes nodepools create"
title: "kubernetes nodepools create"
slug: tcloud_kubernetes_nodepools_create
url: /docs/tcloud/kubernetes/nodepools_create/
weight: 9966
cascade:
  type: docs
---
## tcloud kubernetes nodepools create

Create a Kubernetes node pool

### Synopsis

Create a new node pool in a Kubernetes cluster.

Examples:
  # Create a node pool with minimal configuration
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a

  # Create a node pool with autoscaling
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --enable-autoscaling --min-nodes 1 --max-nodes 5 --availability-zone nl-01a

  # Create node pools in multiple availability zones
  tcloud kubernetes nodepools create --cluster my-cluster --name worker --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a --availability-zone nl-01b

```
tcloud kubernetes nodepools create [flags]
```

### Options

```
      --availability-zone strings   Availability zone for the node pool (can be specified multiple times). If not specified, a random AZ from the cluster's region will be selected.
      --cluster string              Cluster identity, name, or slug (required)
      --enable-autohealing          Enable autohealing for the node pool
      --enable-autoscaling          Enable autoscaling for the node pool
  -h, --help                        help for create
      --machine-type string         Machine type for the node pool (required)
      --max-nodes int               Maximum number of nodes (required when autoscaling is enabled) (default 3)
      --min-nodes int               Minimum number of nodes (required when autoscaling is enabled)
      --name string                 Name of the node pool (default: worker) (default "worker")
      --node-annotations strings    Node annotations in key=value format (applied to Kubernetes nodes)
      --node-labels strings         Node labels in key=value format (applied to Kubernetes nodes)
      --node-taints strings         Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule')
      --num-nodes int               Number of nodes in the node pool (ignored if --enable-autoscaling is set) (default 1)
      --security-groups strings     Security group identities to attach to node pool machines
      --subnet string               Subnet for the node pool (defaults to cluster subnet)
      --upgrade-strategy string     Upgrade strategy: manual, auto, always, on-delete, inplace, or never (default "auto")
      --wait                        Wait for the node pool to be ready before returning
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

