---
linkTitle: "tcloud secrets put"
title: "secrets put"
slug: tcloud_secrets_put
url: /docs/tcloud/secrets/put/
weight: 9696
cascade:
  type: docs
---
## tcloud secrets put

Put a new secret version

```
tcloud secrets put [flags]
```

### Options

```
      --from-file string          Read secret string from a file
      --generate-bytes int[=32]   Generate a random secret of this many bytes (16-4096; bare --generate-bytes uses 32)
  -h, --help                      help for put
      --kv strings                Secret key/value pairs as key=value (repeatable)
      --path string               Secret path
      --region string             Region
      --string string             Secret string value
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

