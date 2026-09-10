---
linkTitle: "tcloud kms decrypt"
title: "kms decrypt"
slug: tcloud_kms_decrypt
url: /docs/tcloud/kms/decrypt/
weight: 9869
cascade:
  type: docs
---
## tcloud kms decrypt

Decrypt ciphertext with a KMS key

### Synopsis

Decrypt data with a KMS key.

Provide exactly one of --ciphertext or --from-file (file containing ciphertext).
When --to-file is set, decoded plaintext bytes are written with mode 0600.
Otherwise base64-encoded plaintext is printed to stdout.

```
tcloud kms decrypt [flags]
```

### Examples

```
  tcloud kms decrypt --region nl-ams --key kms-123 --ciphertext 'thalassa:v1:...'
  tcloud kms decrypt --region nl-ams --key kms-123 --from-file secret.enc --to-file secret.txt
```

### Options

```
      --ciphertext string   Ciphertext from encrypt
      --from-file string    Read ciphertext from a file
  -h, --help                help for decrypt
      --key string          KMS key identity
      --region string       Region
      --to-file string      Write decoded plaintext to a file (mode 0600) instead of stdout
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

* [tcloud kms](/docs/tcloud/tcloud_kms/)	 - Manage KMS keys and cryptographic operations (beta)

