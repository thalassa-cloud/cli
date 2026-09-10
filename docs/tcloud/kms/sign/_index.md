---
linkTitle: "tcloud kms sign"
title: "kms sign"
slug: tcloud_kms_sign
url: /docs/tcloud/kms/sign/
weight: 9854
cascade:
  type: docs
---
## tcloud kms sign

Sign input with an asymmetric KMS key

### Synopsis

Sign data with an asymmetric KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes).
The signature is written to --to-file when set, otherwise to stdout.

```
tcloud kms sign [flags]
```

### Examples

```
  tcloud kms sign --region nl-ams --key kms-123 --input "$(echo -n hello | base64)"
  tcloud kms sign --region nl-ams --key kms-123 --from-file message.txt --to-file message.sig
```

### Options

```
      --context string          Optional signing context
      --from-file string        Read raw input bytes from a file
      --hash-algorithm string   Hash algorithm
  -h, --help                    help for sign
      --input string            Base64-encoded input to sign
      --key string              KMS key identity
      --key-version string      Key version
      --prehashed               Input is already hashed
      --region string           Region
      --to-file string          Write signature to a file (mode 0600) instead of stdout
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
      --debug                  Debug mode
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kms](/docs/tcloud/tcloud_kms/)	 - Manage KMS keys and cryptographic operations (beta)

