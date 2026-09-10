---
linkTitle: "tcloud secrets"
title: "secrets"
slug: tcloud_secrets
url: /docs/tcloud/tcloud_secrets/
weight: 9694
cascade:
  type: docs
---
## tcloud secrets

Manage secrets (beta)

### Synopsis

Manage Secrets Manager paths, versions, and access policies.

Note: This command is in beta. Commands that list or view secrets show
metadata only; use get-value when you intentionally need secret material.

### Options

```
  -h, --help   help for secrets
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud secrets browse](/docs/tcloud/secrets/browse/)	 - Browse secret prefixes and secrets at a path
* [tcloud secrets create](/docs/tcloud/secrets/create/)	 - Create a secret (metadata response only; use get-value to read material)
* [tcloud secrets delete](/docs/tcloud/secrets/delete/)	 - Delete a secret
* [tcloud secrets destroy-version](/docs/tcloud/secrets/destroy-version/)	 - Permanently destroy a secret version
* [tcloud secrets get-value](/docs/tcloud/secrets/get-value/)	 - Print secret material to stdout
* [tcloud secrets list](/docs/tcloud/secrets/list/)	 - List secrets under a path prefix (metadata only)
* [tcloud secrets policy](/docs/tcloud/secrets/policy/)	 - Replace a secret access policy from a JSON file
* [tcloud secrets put](/docs/tcloud/secrets/put/)	 - Put a new secret version
* [tcloud secrets view](/docs/tcloud/secrets/view/)	 - View secret metadata (does not print secret material)

