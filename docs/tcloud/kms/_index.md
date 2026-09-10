---
linkTitle: "tcloud kms"
title: "kms"
slug: tcloud_kms
url: /docs/tcloud/tcloud_kms/
weight: 9849
cascade:
  type: docs
---
## tcloud kms

Manage KMS keys and cryptographic operations (beta)

### Synopsis

Manage Key Management Service keys and cryptographic operations.

Note: This command is in beta.

### Options

```
  -h, --help   help for kms
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
* [tcloud kms decrypt](/docs/tcloud/kms/decrypt/)	 - Decrypt ciphertext with a KMS key
* [tcloud kms encrypt](/docs/tcloud/kms/encrypt/)	 - Encrypt plaintext with a KMS key
* [tcloud kms export](/docs/tcloud/kms/export/)	 - Export key material for a KMS key (when export is allowed)
* [tcloud kms hmac](/docs/tcloud/kms/hmac/)	 - Compute an HMAC with a KMS key
* [tcloud kms keys](/docs/tcloud/kms/keys/)	 - Manage KMS keys
* [tcloud kms public-key](/docs/tcloud/kms/public-key/)	 - Get the public key material for an asymmetric KMS key
* [tcloud kms sign](/docs/tcloud/kms/sign/)	 - Sign input with an asymmetric KMS key
* [tcloud kms summary](/docs/tcloud/kms/summary/)	 - Show KMS availability and regional key counts
* [tcloud kms verify](/docs/tcloud/kms/verify/)	 - Verify a signature with an asymmetric KMS key
* [tcloud kms verify-hmac](/docs/tcloud/kms/verify-hmac/)	 - Verify an HMAC with a KMS key
* [tcloud kms wrapping-key](/docs/tcloud/kms/wrapping-key/)	 - Get the regional wrapping public key for BYOK import

