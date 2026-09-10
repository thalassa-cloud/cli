---
linkTitle: "tcloud kms verify"
title: "kms verify"
slug: tcloud_kms_verify
url: /docs/tcloud/kms/verify/
weight: 9852
cascade:
  type: docs
---
## tcloud kms verify

Verify a signature with an asymmetric KMS key

### Synopsis

Verify a signature with an asymmetric KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes),
and exactly one of --signature or --signature-file. Validity is written to
--to-file when set, otherwise to stdout.

```
tcloud kms verify [flags]
```

### Examples

```
  tcloud kms verify --region nl-ams --key kms-123 --input "$(echo -n hello | base64)" --signature '...'
  tcloud kms verify --region nl-ams --key kms-123 --from-file message.txt --signature-file message.sig
```

### Options

```
      --from-file string        Read raw input bytes from a file
      --hash-algorithm string   Hash algorithm
  -h, --help                    help for verify
      --input string            Base64-encoded input that was signed
      --key string              KMS key identity
      --region string           Region
      --signature string        Signature to verify
      --signature-file string   Read signature from a file
      --to-file string          Write validity (true/false) to a file instead of stdout
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

