---
linkTitle: "tcloud kubernetes connect"
title: "kubernetes connect"
slug: tcloud_kubernetes_connect
url: /docs/tcloud/kubernetes/connect/
weight: 9831
cascade:
  type: docs
---
## tcloud kubernetes connect

Connect your shell to the Kubernetes Cluster

### Synopsis

Connect your shell to a Kubernetes cluster. The cluster argument may be omitted when kubernetes.cluster is set in a directory .thalassa file, or when TCLOUD_CLUSTER_* is already set by a previous connect.

```
tcloud kubernetes connect [cluster] [flags]
```

### Options

```
  -h, --help                     help for connect
      --inline-token             embed the session token in the kubeconfig instead of using a kubectl exec credential plugin
      --kubeconfig-path string   path to write the kubeconfig when --temp=false
      --temp                     use a temporary kubeconfig file that is removed when the shell exits (default true)
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

* [tcloud kubernetes](/docs/tcloud/tcloud_kubernetes/)	 - Manage Kubernetes clusters, node pools and more services related to Kubernetes

