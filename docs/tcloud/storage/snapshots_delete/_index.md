---
linkTitle: "tcloud storage snapshots delete"
title: "storage snapshots delete"
slug: tcloud_storage_snapshots_delete
url: /docs/tcloud/storage/snapshots_delete/
weight: 9939
cascade:
  type: docs
---
## tcloud storage snapshots delete

Delete snapshot(s)

### Synopsis

Delete snapshot(s) by identity or label selector.

```
tcloud storage snapshots delete [flags]
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter snapshots (format: key1=value1,key2=value2)
      --wait              Wait for the snapshot(s) to be deleted
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

* [tcloud storage snapshots](/docs/tcloud/storage/snapshots/)	 - Manage volume snapshots

