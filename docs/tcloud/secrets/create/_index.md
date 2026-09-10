---
linkTitle: "tcloud secrets create"
title: "secrets create"
slug: tcloud_secrets_create
url: /docs/tcloud/secrets/create/
weight: 9702
cascade:
  type: docs
---
## tcloud secrets create

Create a secret (metadata response only; use get-value to read material)

```
tcloud secrets create [flags]
```

### Examples

```
  tcloud secrets create --region nl-01 --path /app/prod/db --kms-key kms-123 --generate-bytes
  tcloud secrets create --region nl-01 --path /app/prod/db --kms-key kms-123 --generate-bytes=64
  tcloud secrets create --region nl-01 --path /app/prod/token --kms-key kms-123 --from-file ./token.txt
```

### Options

```
      --annotations strings       Annotations as key=value (repeatable)
      --description string        Description
      --from-file string          Read secret string from a file
      --generate-bytes int[=32]   Generate a random secret of this many bytes (16-4096; bare --generate-bytes uses 32)
  -h, --help                      help for create
      --kms-key string            KMS key identity used to encrypt the secret
      --kv strings                Secret key/value pairs as key=value (repeatable)
      --labels strings            Labels as key=value (repeatable)
      --no-header                 Do not print table headers
      --path string               Secret path
      --policy-file string        JSON file with an access policy
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

