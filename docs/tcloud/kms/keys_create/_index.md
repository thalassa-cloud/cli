---
linkTitle: "tcloud kms keys create"
title: "kms keys create"
slug: tcloud_kms_keys_create
url: /docs/tcloud/kms/keys_create/
weight: 9864
cascade:
  type: docs
---
## tcloud kms keys create

Create a KMS key

```
tcloud kms keys create [flags]
```

### Options

```
      --allow-rotation               Allow rotation for imported keys
      --annotations strings          Annotations as key=value (repeatable)
      --description string           Key description
      --exact-time                   Show full timestamps instead of relative time
      --export-allowed               Allow exporting key material
      --hash-function string         Hash function for imported or HMAC keys
  -h, --help                         help for create
      --import-key-material string   Wrapped key material for BYOK import
      --key-type string              Key type (aes128-gcm96, aes256-gcm96, chacha20-poly1305, ed25519, ecdsa-p256/384/521, rsa-2048/3072/4096, hmac, hmac-sha256, hmac-sha512)
      --labels strings               Labels as key=value (repeatable)
      --name string                  Key name
      --no-header                    Do not print table headers
      --region string                Region
      --rotation-enabled             Enable automatic key rotation
      --rotation-period-days int     Automatic rotation period in days
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

* [tcloud kms keys](/docs/tcloud/kms/keys/)	 - Manage KMS keys

