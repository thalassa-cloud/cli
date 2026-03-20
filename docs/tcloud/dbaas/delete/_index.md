---
linkTitle: "tcloud dbaas delete"
title: "dbaas delete"
slug: tcloud_dbaas_delete
url: /docs/tcloud/dbaas/delete/
weight: 9963
cascade:
  type: docs
---
## tcloud dbaas delete

Delete database cluster(s)

### Synopsis

Delete database cluster(s) by identity or label selector.

```
tcloud dbaas delete [flags]
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter clusters (format: key1=value1,key2=value2)
      --wait              Wait for the database cluster(s) to be deleted
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

* [tcloud dbaas](/docs/tcloud/tcloud_dbaas/)	 - Manage database clusters and related services

