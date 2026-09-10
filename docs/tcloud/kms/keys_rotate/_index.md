---
linkTitle: "tcloud kms keys rotate"
title: "kms keys rotate"
slug: tcloud_kms_keys_rotate
url: /docs/tcloud/kms/keys_rotate/
weight: 9859
cascade:
  type: docs
---
## tcloud kms keys rotate

Rotate a KMS key on demand

```
tcloud kms keys rotate <key> [flags]
```

### Options

```
      --exact-time      Show full timestamps instead of relative time
  -h, --help            help for rotate
      --no-header       Do not print table headers
      --region string   Region
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

