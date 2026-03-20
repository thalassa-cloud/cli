---
linkTitle: "tcloud iam roles rules delete"
title: "iam roles rules delete"
slug: tcloud_iam_roles_rules_delete
url: /docs/tcloud/iam/roles_rules_delete/
weight: 9929
cascade:
  type: docs
---
## tcloud iam roles rules delete

Remove a permission rule from a role

```
tcloud iam roles rules delete <role> <rule> [flags]
```

### Options

```
      --force       Skip the confirmation prompt and delete
  -h, --help        help for delete
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam roles rules](/docs/tcloud/iam/roles_rules/)	 - Permission rules on a role

