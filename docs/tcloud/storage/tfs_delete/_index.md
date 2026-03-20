---
linkTitle: "tcloud storage tfs delete"
title: "storage tfs delete"
slug: tcloud_storage_tfs_delete
url: /docs/tcloud/storage/tfs_delete/
weight: 9843
cascade:
  type: docs
---
## tcloud storage tfs delete

Delete TFS instance(s)

### Synopsis

Delete TFS instance(s) by identity or label selector.

```
tcloud storage tfs delete [flags]
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter TFS instances (format: key1=value1,key2=value2)
      --wait              Wait for the TFS instance(s) to be deleted
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

* [tcloud storage tfs](/docs/tcloud/storage/tfs/)	 - Manage TFS (Thalassa File System) instances

