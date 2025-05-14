---
date: 2025-05-14T17:58:13+02:00
linkTitle: "tcloud kubernetes upgrade"
title: "kubernetes upgrade"
slug: tcloud_kubernetes_upgrade
url: /docs/tcloud/tcloud_kubernetes_upgrade/
weight: 9977
---
## tcloud kubernetes upgrade

Upgrade a Kubernetes cluster

```
tcloud kubernetes upgrade <cluster> [flags]
```

### Options

```
      --all              upgrade all clusters
      --dry-run          print the upgrade request without actually upgrading the cluster (default true)
  -h, --help             help for upgrade
  -v, --version string   the version to upgrade to. If not provided, the latest suitable version will be used following Kubernetes version policy (only +1 minor version or patch updates)
```

### Options inherited from parent commands

```
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes](/docs/tcloud/tcloud_kubernetes/)	 - Manage Kubernetes clusters, node pools and more services related to Kubernetes

