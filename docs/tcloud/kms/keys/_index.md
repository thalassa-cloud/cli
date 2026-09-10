---
linkTitle: "tcloud kms keys"
title: "kms keys"
slug: tcloud_kms_keys
url: /docs/tcloud/kms/keys/
weight: 9856
cascade:
  type: docs
---
## tcloud kms keys

Manage KMS keys

### Options

```
  -h, --help   help for keys
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
* [tcloud kms keys cancel-deletion](/docs/tcloud/kms/keys_cancel-deletion/)	 - Cancel a pending KMS key deletion
* [tcloud kms keys create](/docs/tcloud/kms/keys_create/)	 - Create a KMS key
* [tcloud kms keys delete](/docs/tcloud/kms/keys_delete/)	 - Schedule a KMS key for deletion
* [tcloud kms keys disable](/docs/tcloud/kms/keys_disable/)	 - Disable a KMS key
* [tcloud kms keys enable](/docs/tcloud/kms/keys_enable/)	 - Enable a KMS key
* [tcloud kms keys list](/docs/tcloud/kms/keys_list/)	 - List KMS keys in a region
* [tcloud kms keys rotate](/docs/tcloud/kms/keys_rotate/)	 - Rotate a KMS key on demand
* [tcloud kms keys rotation](/docs/tcloud/kms/keys_rotation/)	 - Update automatic rotation settings for a KMS key
* [tcloud kms keys view](/docs/tcloud/kms/keys_view/)	 - View a KMS key

