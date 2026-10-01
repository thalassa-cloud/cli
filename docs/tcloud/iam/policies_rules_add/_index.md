---
linkTitle: "tcloud iam policies rules add"
title: "iam policies rules add"
slug: tcloud_iam_policies_rules_add
url: /docs/tcloud/iam/policies_rules_add/
weight: 9891
cascade:
  type: docs
---
## tcloud iam policies rules add

Add a permission rule to a policy

```
tcloud iam policies rules add <policy> [flags]
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
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam policies rules](/docs/tcloud/iam/policies_rules/)	 - Permission rules on a policy

