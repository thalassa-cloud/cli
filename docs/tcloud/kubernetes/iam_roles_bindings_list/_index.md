---
linkTitle: "tcloud kubernetes iam roles bindings list"
title: "kubernetes iam roles bindings list"
slug: tcloud_kubernetes_iam_roles_bindings_list
url: /docs/tcloud/kubernetes/iam_roles_bindings_list/
weight: 9843
cascade:
  type: docs
---
## tcloud kubernetes iam roles bindings list

List bindings for a Kubernetes cluster role

```
tcloud kubernetes iam roles bindings list <role> [flags]
```

### Options

```
  -h, --help        help for list
      --no-header   Do not print table headers
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

* [tcloud kubernetes iam roles bindings](/docs/tcloud/kubernetes/iam_roles_bindings/)	 - Kubernetes cluster role bindings (who receives the role)

