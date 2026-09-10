---
linkTitle: "tcloud kms public-key"
title: "kms public-key"
slug: tcloud_kms_public-key
url: /docs/tcloud/kms/public-key/
weight: 9855
cascade:
  type: docs
---
## tcloud kms public-key

Get the public key material for an asymmetric KMS key

```
tcloud kms public-key [flags]
```

### Options

```
  -h, --help            help for public-key
      --key string      KMS key identity
      --no-header       Do not print table headers
      --region string   Region
      --version int     Key version
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

