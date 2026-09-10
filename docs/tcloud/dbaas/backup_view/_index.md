---
linkTitle: "tcloud dbaas backup view"
title: "dbaas backup view"
slug: tcloud_dbaas_backup_view
url: /docs/tcloud/dbaas/backup_view/
weight: 9969
cascade:
  type: docs
---
## tcloud dbaas backup view

View backup details

### Synopsis

View detailed information about a database backup

```
tcloud dbaas backup view [flags]
```

### Options

```
      --exact-time   Show exact time instead of relative time
  -h, --help         help for view
      --no-header    Do not print the header
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

* [tcloud dbaas backup](/docs/tcloud/dbaas/backup/)	 - Manage database backups

