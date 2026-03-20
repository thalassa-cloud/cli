---
linkTitle: "tcloud dbaas backup create"
title: "dbaas backup create"
slug: tcloud_dbaas_backup_create
url: /docs/tcloud/dbaas/backup_create/
weight: 9975
cascade:
  type: docs
---
## tcloud dbaas backup create

Create a database backup

### Synopsis

Create a new backup for a database cluster

```
tcloud dbaas backup create [flags]
```

### Options

```
      --annotations strings       Annotations in key=value format (can be specified multiple times)
      --description string        Description of the backup
  -h, --help                      help for create
      --labels strings            Labels in key=value format (can be specified multiple times)
      --name string               Name of the backup (required)
      --no-header                 Do not print the header
      --retention-policy string   Retention policy for the backup
      --wait                      Wait for the backup to be completed before returning
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

* [tcloud dbaas backup](/docs/tcloud/dbaas/backup/)	 - Manage database backups

