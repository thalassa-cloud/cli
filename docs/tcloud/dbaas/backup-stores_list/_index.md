---
linkTitle: "tcloud dbaas backup-stores list"
title: "dbaas backup-stores list"
slug: tcloud_dbaas_backup-stores_list
url: /docs/tcloud/dbaas/backup-stores_list/
weight: 9959
cascade:
  type: docs
---
## tcloud dbaas backup-stores list

List database backup stores

### Synopsis

List backup stores in the organisation

```
tcloud dbaas backup-stores list [flags]
```

### Options

```
      --exact-time    Show exact time instead of relative time
  -h, --help          help for list
      --no-header     Do not print the header
      --show-labels   Show labels
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

* [tcloud dbaas backup-stores](/docs/tcloud/dbaas/backup-stores/)	 - Manage database backup stores

