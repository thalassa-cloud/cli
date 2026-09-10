---
linkTitle: "tcloud kubernetes iam roles create"
title: "kubernetes iam roles create"
slug: tcloud_kubernetes_iam_roles_create
url: /docs/tcloud/kubernetes/iam_roles_create/
weight: 9841
cascade:
  type: docs
---
## tcloud kubernetes iam roles create

Create a custom Kubernetes cluster IAM role

```
tcloud kubernetes iam roles create [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Role description
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --name string           Role name
      --no-header             Do not print table headers
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes iam roles](/docs/tcloud/kubernetes/iam_roles/)	 - Kubernetes cluster IAM roles, permission rules, and bindings

