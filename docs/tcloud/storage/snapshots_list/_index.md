---
linkTitle: "tcloud storage snapshots list"
title: "storage snapshots list"
slug: tcloud_storage_snapshots_list
url: /docs/tcloud/storage/snapshots_list/
weight: 9928
cascade:
  type: docs
---
## tcloud storage snapshots list

Get a list of snapshots

```
tcloud storage snapshots list [flags]
```

### Options

```
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter snapshots (format: key1=value1,key2=value2)
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

* [tcloud storage snapshots](/docs/tcloud/storage/snapshots/)	 - Manage volume snapshots

