---
linkTitle: "tcloud kubernetes nodepools upgrade"
title: "kubernetes nodepools upgrade"
slug: tcloud_kubernetes_nodepools_upgrade
url: /docs/tcloud/kubernetes/nodepools_upgrade/
weight: 9809
cascade:
  type: docs
---
## tcloud kubernetes nodepools upgrade

Upgrade a Kubernetes node pool version

### Synopsis

Upgrade the Kubernetes version of a node pool.

The node pool update API replaces the whole node pool. This command changes only the Kubernetes version and resends the current machine type, replicas, autoscaling, node settings, and security groups.

The target version follows the Kubernetes version policy: the newest patch of the current minor, otherwise the newest patch of the next minor. A node pool is not upgraded past the cluster control plane version.

Examples:
  # Show the version a node pool would move to
  tcloud kubernetes nodepools upgrade --cluster my-cluster --name worker

  # Upgrade a node pool to a specific version
  tcloud kubernetes nodepools upgrade --cluster my-cluster --name worker --version 1.32.2 --dry-run=false

  # Upgrade every node pool in a cluster
  tcloud kubernetes nodepools upgrade --cluster my-cluster --all --dry-run=false

```
tcloud kubernetes nodepools upgrade [flags]
```

### Options

```
      --all                     Upgrade all node pools in the cluster
      --cluster string          Cluster identity, name, or slug (defaults to kubernetes.cluster in .thalassa)
      --dry-run                 Print the upgrade without changing the node pool (default true)
  -h, --help                    help for upgrade
      --name string             Node pool name, identity, or slug (required unless --all is set)
  -v, --version string          Version to upgrade to (identity, slug, name, or Kubernetes version). Defaults to the newest patch, or the next minor when no patch is available
      --wait                    Wait for the node pool upgrade to complete
      --wait-timeout duration   Maximum time to wait for the node pool upgrade to complete (default 20m0s)
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
