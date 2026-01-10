---
linkTitle: "tcloud kubernetes update"
title: "kubernetes update"
slug: tcloud_kubernetes_update
url: /docs/tcloud/kubernetes/update/
weight: 9960
cascade:
  type: docs
---
## tcloud kubernetes update

Update a Kubernetes cluster

### Synopsis

Update an existing Kubernetes cluster.

Only fields specified with flags will be updated. All other fields will remain unchanged.

Examples:
  # Update cluster description
  tcloud kubernetes update my-cluster --description "Production cluster"

  # Update Kubernetes version
  tcloud kubernetes update my-cluster --cluster-version 1.28.0

  # Update maintenance window
  tcloud kubernetes update my-cluster --maintenance-day monday --maintenance-start 02:00

  # Update multiple fields
  tcloud kubernetes update my-cluster --description "Updated description" --disable-public-endpoint

Note: Use 'tcloud kubernetes label' and 'tcloud kubernetes annotate' commands
to manage labels and annotations separately.

```
tcloud kubernetes update <cluster> [flags]
```

### Options

```
      --audit-log-profile string        Audit log profile: none, basic, or metadata
      --cluster-version string          Kubernetes version (name, slug, or identity)
      --default-network-policy string   Default network policy: allow-all, deny-all, or none
      --description string              Description of the cluster
      --disable-public-endpoint         Disable public API server endpoint
  -h, --help                            help for update
      --kube-proxy-deployment string    Kube proxy deployment: disabled, managed, or custom
      --kube-proxy-mode string          Kube proxy mode: iptables or ipvs
      --maintenance-day string          Maintenance day: 0-6, Sunday-Saturday, or day name
      --maintenance-start string        Maintenance start time: HH:MM format (e.g., '02:00' or '14:30')
      --name string                     Name of the cluster
      --pod-security-standards string   Pod security standards profile: baseline, restricted, or privileged
      --wait                            Wait for the cluster update to complete
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

