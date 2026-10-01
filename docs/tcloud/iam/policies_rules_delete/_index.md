---
linkTitle: "tcloud iam policies rules delete"
title: "iam policies rules delete"
slug: tcloud_iam_policies_rules_delete
url: /docs/tcloud/iam/policies_rules_delete/
weight: 9890
cascade:
  type: docs
---
## tcloud iam policies rules delete

Remove a permission rule from a policy

```
tcloud iam policies rules delete <policy> <rule> [flags]
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
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam policies rules](/docs/tcloud/iam/policies_rules/)	 - Permission rules on a policy

