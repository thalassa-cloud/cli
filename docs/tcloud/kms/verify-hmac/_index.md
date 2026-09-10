---
linkTitle: "tcloud kms verify-hmac"
title: "kms verify-hmac"
slug: tcloud_kms_verify-hmac
url: /docs/tcloud/kms/verify-hmac/
weight: 9851
cascade:
  type: docs
---
## tcloud kms verify-hmac

Verify an HMAC with a KMS key

### Synopsis

Verify an HMAC with a KMS key.

Provide exactly one of --input (base64-encoded) or --from-file (raw file bytes),
and exactly one of --hmac or --hmac-file. Validity is written to --to-file when
set, otherwise to stdout.

```
tcloud kms verify-hmac [flags]
```

### Examples

```
  tcloud kms verify-hmac --region nl-ams --key kms-123 --from-file message.txt --hmac-file message.hmac
```

### Options

```
      --from-file string        Read raw input bytes from a file
      --hash-algorithm string   Hash algorithm
  -h, --help                    help for verify-hmac
      --hmac string             HMAC value to verify
      --hmac-file string        Read HMAC value from a file
      --input string            Base64-encoded input that was HMACed
      --key string              KMS key identity
      --region string           Region
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

