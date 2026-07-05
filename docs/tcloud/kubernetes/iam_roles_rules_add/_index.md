---
linkTitle: "tcloud kubernetes iam roles rules add"
title: "kubernetes iam roles rules add"
slug: tcloud_kubernetes_iam_roles_rules_add
url: /docs/tcloud/kubernetes/iam_roles_rules_add/
weight: 9891
cascade:
  type: docs
---
## tcloud kubernetes iam roles rules add

Add a permission rule to a Kubernetes cluster role

```
tcloud kubernetes iam roles rules add <role> [flags]
```

### Options

```
      --api-group strings          API group (repeatable)
  -h, --help                       help for add
      --no-header                  Do not print table headers
      --non-resource-url strings   Non-resource URL (repeatable)
      --note string                Human-readable note for the rule
      --resource strings           Resource name (repeatable)
      --resource-name strings      Concrete resource name (repeatable)
      --verb strings               Verb: get, list, watch, create, update, delete, patch, or * (repeatable)
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

