---
linkTitle: "tcloud storage snapshot-policies list"
title: "storage snapshot-policies list"
slug: tcloud_storage_snapshot-policies_list
url: /docs/tcloud/storage/snapshot-policies_list/
weight: 9691
cascade:
  type: docs
---
## tcloud storage snapshot-policies list

List snapshot policies

```
tcloud storage snapshot-policies list [flags]
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Filter by region
  -l, --selector string   Label selector to filter snapshot policies
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud storage snapshot-policies](/docs/tcloud/storage/snapshot-policies/)	 - Manage snapshot policies

