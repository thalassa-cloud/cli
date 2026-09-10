---
linkTitle: "tcloud kms keys delete"
title: "kms keys delete"
slug: tcloud_kms_keys_delete
url: /docs/tcloud/kms/keys_delete/
weight: 9863
cascade:
  type: docs
---
## tcloud kms keys delete

Schedule a KMS key for deletion

```
tcloud kms keys delete <key> [flags]
```

### Options

```
      --force           Skip the confirmation prompt and delete
  -h, --help            help for delete
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

