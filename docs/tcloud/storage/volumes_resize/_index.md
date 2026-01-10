---
linkTitle: "tcloud storage volumes resize"
title: "storage volumes resize"
slug: tcloud_storage_volumes_resize
url: /docs/tcloud/storage/volumes_resize/
weight: 9909
cascade:
  type: docs
---
## tcloud storage volumes resize

Resize volume(s)

### Synopsis

Resize volume(s) to a new size in GB. The new size must be larger than the current size. Can resize multiple volumes by identity or using a label selector.

```
tcloud storage volumes resize [volume-id...] [flags]
```

### Options

```
      --force             Force the resize and skip the confirmation
  -h, --help              help for resize
  -l, --selector string   Label selector to filter volumes (format: key1=value1,key2=value2)
      --size int          New size in GB (required)
      --wait              Wait for the resize operation to complete
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

