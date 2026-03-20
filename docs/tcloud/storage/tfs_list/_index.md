---
linkTitle: "tcloud storage tfs list"
title: "storage tfs list"
slug: tcloud_storage_tfs_list
url: /docs/tcloud/storage/tfs_list/
weight: 9842
cascade:
  type: docs
---
## tcloud storage tfs list

Get a list of TFS instances

```
tcloud storage tfs list [flags]
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Region of the TFS instance
  -l, --selector string   Label selector to filter TFS instances (format: key1=value1,key2=value2)
      --show-labels       Show labels
      --status string     Status of the TFS instance
      --vpc string        VPC of the TFS instance
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

