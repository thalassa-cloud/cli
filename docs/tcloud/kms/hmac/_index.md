---
linkTitle: "tcloud kms hmac"
title: "kms hmac"
slug: tcloud_kms_hmac
url: /docs/tcloud/kms/hmac/
weight: 9866
cascade:
  type: docs
---
## tcloud kms hmac

Compute an HMAC with a KMS key

```
tcloud kms hmac [flags]
```

### Options

```
      --algorithm string     HMAC algorithm
  -h, --help                 help for hmac
      --input string         Input for HMAC
      --key string           KMS key identity
      --key-version string   Key version
      --region string        Region
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

