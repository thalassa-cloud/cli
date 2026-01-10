---
linkTitle: "tcloud kubernetes create"
title: "kubernetes create"
slug: tcloud_kubernetes_create
url: /docs/tcloud/kubernetes/create/
weight: 9972
cascade:
  type: docs
---
## tcloud kubernetes create

Create a Kubernetes cluster

### Synopsis

Create a Kubernetes cluster in the Thalassa Cloud Platform.

This command creates a new Kubernetes cluster with sensible defaults.

Examples:
  # Create a managed cluster with minimal configuration
  tcloud kubernetes create my-cluster --subnet subnet-123

  # Create a cluster with custom networking
  tcloud kubernetes create my-cluster --subnet subnet-123 --pod-cidr 10.0.0.0/16 --service-cidr 172.16.0.0/18

  # Create a cluster and wait for it to be ready
  tcloud kubernetes create my-cluster --subnet subnet-123 --wait

  # Create a cluster with a node pool
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a

  # Create a cluster with node pools in multiple availability zones
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --num-nodes 3 --availability-zone nl-01a --availability-zone nl-01b

  # Create a cluster with autoscaling node pool (auto-selects AZ if not provided)
  tcloud kubernetes create my-cluster --subnet subnet-123 --machine-type pgp-medium --enable-autoscaling --min-nodes 1 --max-nodes 5

  # Create a cluster with scheduled upgrades (maintenance window)
  tcloud kubernetes create my-cluster --subnet subnet-123 --maintenance-day monday --maintenance-start 02:00

```
tcloud kubernetes create [flags]
```

### Options

```
      --annotations strings             Annotations in key=value format (can be specified multiple times)
      --audit-log-profile string        Audit log profile: none, basic, or metadata (default: none)
      --availability-zone strings       Availability zone for the node pool (can be specified multiple times to create node pools in multiple AZs). If not specified, a random AZ from the cluster's region will be selected.
      --cluster-type string             Cluster type: managed or hosted-control-plane (default "managed")
      --cluster-version string          Kubernetes version (name, slug, or identity). Defaults to latest stable
      --cni string                      CNI plugin: cilium or custom (default "cilium")
      --default-network-policy string   Default network policy: allow-all, deny-all, or none (default: allow-all)
      --description string              Description of the cluster
      --disable-public-endpoint         Disable public API server endpoint
      --enable-autohealing              Enable autohealing for the node pool
      --enable-autoscaling              Enable autoscaling for the node pool
  -h, --help                            help for create
      --kube-proxy-deployment string    Kube proxy deployment: disabled, managed, or custom (default: disabled)
      --kube-proxy-mode string          Kube proxy mode: iptables or ipvs (default: iptables)
      --labels strings                  Labels in key=value format (can be specified multiple times)
      --machine-type string             Machine type for the node pool (required to create a node pool)
      --max-nodes int                   Maximum number of nodes (required when autoscaling is enabled) (default 3)
      --min-nodes int                   Minimum number of nodes (required when autoscaling is enabled) (default 1)
      --name string                     Name of the cluster (deprecated: use positional argument)
      --node-annotations strings        Node annotations in key=value format (applied to Kubernetes nodes)
      --node-labels strings             Node labels in key=value format (applied to Kubernetes nodes)
      --node-pool-name string           Name of the node pool (default "worker")
      --node-pool-subnet string         Subnet for the node pool (defaults to cluster subnet)
      --node-taints strings             Node taints in key=value:effect or key:effect format (e.g., 'dedicated=gpu:NoSchedule')
      --num-nodes int                   Number of nodes in the node pool (ignored if --enable-autoscaling is set) (default 1)
      --pod-cidr string                 Pod CIDR (default "192.168.0.0/16")
      --pod-security-standards string   Pod security standards profile: baseline, restricted, or privileged (default: baseline)
      --region string                   Region for hosted-control-plane clusters
      --security-groups strings         Security group identities to attach to node pool machines
      --service-cidr string             Service CIDR (default "172.16.0.0/18")
      --subnet string                   Subnet for managed clusters
      --upgrade-strategy string         Upgrade strategy: manual, auto, always, on-delete, inplace, or never (default "auto")
      --wait                            Wait for the cluster to be ready before returning
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

