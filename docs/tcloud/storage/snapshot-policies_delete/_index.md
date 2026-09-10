---
linkTitle: "tcloud storage snapshot-policies delete"
title: "storage snapshot-policies delete"
slug: tcloud_storage_snapshot-policies_delete
url: /docs/tcloud/storage/snapshot-policies_delete/
weight: 9692
cascade:
  type: docs
---
## tcloud storage snapshot-policies delete

Delete snapshot policy(ies)

```
tcloud storage snapshot-policies delete [flags]
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter snapshot policies
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

* [tcloud storage snapshot-policies](/docs/tcloud/storage/snapshot-policies/)	 - Manage snapshot policies

