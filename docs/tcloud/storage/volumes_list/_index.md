---
linkTitle: "tcloud storage volumes list"
title: "storage volumes list"
slug: tcloud_storage_volumes_list
url: /docs/tcloud/storage/volumes_list/
weight: 9933
cascade:
  type: docs
---
## tcloud storage volumes list

Get a list of volumes

```
tcloud storage volumes list [flags]
```

### Options

```
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter volumes (format: key1=value1,key2=value2)
      --show-labels       Show labels
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

