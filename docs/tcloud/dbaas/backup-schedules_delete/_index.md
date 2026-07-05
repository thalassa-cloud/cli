---
linkTitle: "tcloud dbaas backup-schedules delete"
title: "dbaas backup-schedules delete"
slug: tcloud_dbaas_backup-schedules_delete
url: /docs/tcloud/dbaas/backup-schedules_delete/
weight: 9967
cascade:
  type: docs
---
## tcloud dbaas backup-schedules delete

Delete a database backup schedule

### Synopsis

Delete a database backup schedule

```
tcloud dbaas backup-schedules delete [flags]
```

### Options

```
      --force   Force the deletion and skip the confirmation
  -h, --help    help for delete
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

* [tcloud dbaas backup-schedules](/docs/tcloud/dbaas/backup-schedules/)	 - Manage database backup schedules

