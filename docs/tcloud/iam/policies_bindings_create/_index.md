---
linkTitle: "tcloud iam policies bindings create"
title: "iam policies bindings create"
slug: tcloud_iam_policies_bindings_create
url: /docs/tcloud/iam/policies_bindings_create/
weight: 9900
cascade:
  type: docs
---
## tcloud iam policies bindings create

Create a binding to a user or service account

```
tcloud iam policies bindings create <policy> [flags]
```

### Options

```
      --annotations strings               Annotations as key=value (repeatable)
      --description string                Binding description
      --duration string                   Relative TTL (e.g. 24h, 7d) instead of --expires-at
      --expires-at string                 Absolute expiry time (RFC3339)
  -h, --help                              help for create
      --labels strings                    Labels as key=value (repeatable)
      --name string                       Binding name
      --no-header                         Do not print table headers
      --service-account-identity string   Bind to this service account identity
      --user-identity string              Bind to this user identity
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

