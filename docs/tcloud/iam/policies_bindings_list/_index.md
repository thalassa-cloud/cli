---
linkTitle: "tcloud iam policies bindings list"
title: "iam policies bindings list"
slug: tcloud_iam_policies_bindings_list
url: /docs/tcloud/iam/policies_bindings_list/
weight: 9898
cascade:
  type: docs
---
## tcloud iam policies bindings list

List bindings for a policy

```
tcloud iam policies bindings list <policy> [flags]
```

### Options

```
      --exact-time   Show full timestamps instead of relative time
  -h, --help         help for list
      --no-header    Do not print table headers
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

