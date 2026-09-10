---
linkTitle: "tcloud kms keys rotation"
title: "kms keys rotation"
slug: tcloud_kms_keys_rotation
url: /docs/tcloud/kms/keys_rotation/
weight: 9858
cascade:
  type: docs
---
## tcloud kms keys rotation

Update automatic rotation settings for a KMS key

```
tcloud kms keys rotation <key> [flags]
```

### Options

```
      --enabled           Enable or disable automatic rotation
      --exact-time        Show full timestamps instead of relative time
  -h, --help              help for rotation
      --no-header         Do not print table headers
      --period-days int   Automatic rotation period in days
      --region string     Region
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

