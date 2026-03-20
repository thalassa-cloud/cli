---
linkTitle: "tcloud dbaas backup-schedules create"
title: "dbaas backup-schedules create"
slug: tcloud_dbaas_backup-schedules_create
url: /docs/tcloud/dbaas/backup-schedules_create/
weight: 9970
cascade:
  type: docs
---
## tcloud dbaas backup-schedules create

Create a database backup schedule

### Synopsis

Create a new backup schedule for a database cluster

```
tcloud dbaas backup-schedules create [flags]
```

### Options

```
      --annotations strings       Annotations in key=value format (can be specified multiple times)
      --description string        Description of the backup schedule
  -h, --help                      help for create
      --labels strings            Labels in key=value format (can be specified multiple times)
      --method string             Backup method: 'barman' (default) (default "barman")
  -n, --name string               Name of the backup schedule (required)
      --no-header                 Do not print the header
      --retention-policy string   Retention policy for the backup schedule (required)
      --schedule string           Cron expression for the backup schedule (required)
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

