---
linkTitle: "tcloud secrets get-value"
title: "secrets get-value"
slug: tcloud_secrets_get-value
url: /docs/tcloud/secrets/get-value/
weight: 9699
cascade:
  type: docs
---
## tcloud secrets get-value

Print secret material to stdout

### Synopsis

Fetches and prints the secret value. Prefer piping to a file and avoid logging the output.

```
tcloud secrets get-value [flags]
```

### Options

```
  -h, --help            help for get-value
      --path string     Secret path
      --region string   Region
      --version int     Specific version (defaults to current)
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

* [tcloud secrets](/docs/tcloud/tcloud_secrets/)	 - Manage secrets (beta)

