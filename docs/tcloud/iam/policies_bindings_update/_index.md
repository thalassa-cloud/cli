---
linkTitle: "tcloud iam policies bindings update"
title: "iam policies bindings update"
slug: tcloud_iam_policies_bindings_update
url: /docs/tcloud/iam/policies_bindings_update/
weight: 9897
cascade:
  type: docs
---
## tcloud iam policies bindings update

Update a policy binding (only set flags are changed)

```
tcloud iam policies bindings update <policy> <binding> [flags]
```

### Options

```
      --annotations strings   Replace annotations (key=value, repeatable)
      --description string    Binding description
  -h, --help                  help for update
      --labels strings        Replace labels (key=value, repeatable)
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
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam policies bindings](/docs/tcloud/iam/policies_bindings/)	 - Policy bindings (who receives the policy)

