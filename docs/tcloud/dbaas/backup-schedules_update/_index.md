---
linkTitle: "tcloud dbaas backup-schedules update"
title: "dbaas backup-schedules update"
slug: tcloud_dbaas_backup-schedules_update
url: /docs/tcloud/dbaas/backup-schedules_update/
weight: 9967
cascade:
  type: docs
---
## tcloud dbaas backup-schedules update

Update a database backup schedule

### Synopsis

Update properties of an existing backup schedule

```
tcloud dbaas backup-schedules update [flags]
```

### Options

```
      --annotations strings       Annotations in key=value format (can be specified multiple times)
      --description string        Description of the backup schedule
  -h, --help                      help for update
      --labels strings            Labels in key=value format (can be specified multiple times)
      --name string               Name of the backup schedule
      --no-header                 Do not print the header
      --retention-policy string   Retention policy for the backup schedule
      --schedule string           Cron expression for the backup schedule
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

* [tcloud dbaas backup-schedules](/docs/tcloud/dbaas/backup-schedules/)	 - Manage database backup schedules

