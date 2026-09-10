---
linkTitle: "tcloud kubernetes iam roles bindings delete"
title: "kubernetes iam roles bindings delete"
slug: tcloud_kubernetes_iam_roles_bindings_delete
url: /docs/tcloud/kubernetes/iam_roles_bindings_delete/
weight: 9844
cascade:
  type: docs
---
## tcloud kubernetes iam roles bindings delete

Delete a Kubernetes cluster role binding

```
tcloud kubernetes iam roles bindings delete <role> <binding> [flags]
```

### Options

```
      --force   Skip the confirmation prompt and delete
  -h, --help    help for delete
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

