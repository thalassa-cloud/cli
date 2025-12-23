---
linkTitle: "tcloud storage volumes delete"
title: "storage volumes delete"
slug: tcloud_storage_volumes_delete
url: /docs/tcloud/storage/volumes_delete/
weight: 9935
cascade:
  type: docs
---
## tcloud storage volumes delete

Delete volume(s)

### Synopsis

Delete volume(s) by identity or label selector.

```
tcloud storage volumes delete [flags]
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter volumes (format: key1=value1,key2=value2)
      --wait              Wait for the volume(s) to be deleted
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud storage volumes](/docs/tcloud/storage/volumes/)	 - Manage storage volumes

