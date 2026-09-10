---
linkTitle: "tcloud dbaas backup-stores create"
title: "dbaas backup-stores create"
slug: tcloud_dbaas_backup-stores_create
url: /docs/tcloud/dbaas/backup-stores_create/
weight: 9961
cascade:
  type: docs
---
## tcloud dbaas backup-stores create

Create a database backup store

### Synopsis

Create a backup store for Barman backups and point-in-time recovery

```
tcloud dbaas backup-stores create [flags]
```

### Options

```
      --annotations strings       Annotations in key=value format (can be specified multiple times)
      --delete-protection         Enable delete protection
      --description string        Description of the backup store
  -h, --help                      help for create
      --labels strings            Labels in key=value format (can be specified multiple times)
      --name string               Name of the backup store (required)
      --no-header                 Do not print the header
      --region string             Region identity or slug (required)
      --retention-mode string     Retention mode: retainForPointInTime or forceCleanupAfterExpiry (default "retainForPointInTime")
      --retention-policy string   Retention policy in days, for example 30d
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

