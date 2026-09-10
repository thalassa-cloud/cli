---
linkTitle: "tcloud kms encrypt"
title: "kms encrypt"
slug: tcloud_kms_encrypt
url: /docs/tcloud/kms/encrypt/
weight: 9868
cascade:
  type: docs
---
## tcloud kms encrypt

Encrypt plaintext with a KMS key

### Synopsis

Encrypt data with a KMS key.

Provide exactly one of --plaintext (base64-encoded) or --from-file (raw file
bytes). Ciphertext is written to --to-file when set, otherwise to stdout.

```
tcloud kms encrypt [flags]
```

### Examples

```
  tcloud kms encrypt --region nl-ams --key kms-123 --plaintext "$(echo -n hello | base64)"
  tcloud kms encrypt --region nl-ams --key kms-123 --from-file secret.txt --to-file secret.enc
```

### Options

```
      --from-file string     Read raw plaintext bytes from a file
  -h, --help                 help for encrypt
      --key string           KMS key identity
      --key-version string   Key version
      --plaintext string     Base64-encoded plaintext
      --region string        Region
      --to-file string       Write ciphertext to a file (mode 0600) instead of stdout
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

