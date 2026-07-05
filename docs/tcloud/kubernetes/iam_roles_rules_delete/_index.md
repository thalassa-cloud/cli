---
linkTitle: "tcloud kubernetes iam roles rules delete"
title: "kubernetes iam roles rules delete"
slug: tcloud_kubernetes_iam_roles_rules_delete
url: /docs/tcloud/kubernetes/iam_roles_rules_delete/
weight: 9890
cascade:
  type: docs
---
## tcloud kubernetes iam roles rules delete

Remove a permission rule from a Kubernetes cluster role

```
tcloud kubernetes iam roles rules delete <role> <rule> [flags]
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

* [tcloud kubernetes iam roles rules](/docs/tcloud/kubernetes/iam_roles_rules/)	 - Permission rules on a Kubernetes cluster role

