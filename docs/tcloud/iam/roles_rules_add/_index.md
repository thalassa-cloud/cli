---
linkTitle: "tcloud iam roles rules add"
title: "iam roles rules add"
slug: tcloud_iam_roles_rules_add
url: /docs/tcloud/iam/roles_rules_add/
weight: 9930
cascade:
  type: docs
---
## tcloud iam roles rules add

Add a permission rule to a role

```
tcloud iam roles rules add <role> [flags]
```

### Options

```
  -h, --help                        help for add
      --no-header                   Do not print table headers
      --note string                 Human-readable note for the rule
      --permission strings          Permission: create, read, update, delete, list, or * (repeatable)
      --resource strings            Resource type (repeatable)
      --resource-identity strings   Concrete resource identity (repeatable)
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

* [tcloud iam roles rules](/docs/tcloud/iam/roles_rules/)	 - Permission rules on a role

