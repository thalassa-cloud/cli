---
linkTitle: "tcloud dbaas backup-stores update"
title: "dbaas backup-stores update"
slug: tcloud_dbaas_backup-stores_update
url: /docs/tcloud/dbaas/backup-stores_update/
weight: 9958
cascade:
  type: docs
---
## tcloud dbaas backup-stores update

Update a database backup store

### Synopsis

Update properties of an existing backup store

```
tcloud dbaas backup-stores update [flags]
```

### Options

```
      --annotations strings       Annotations in key=value format (can be specified multiple times)
      --delete-protection         Enable or disable delete protection
      --description string        Description of the backup store
  -h, --help                      help for update
      --labels strings            Labels in key=value format (can be specified multiple times)
      --name string               Name of the backup store
      --no-header                 Do not print the header
      --retention-mode string     Retention mode: retainForPointInTime or forceCleanupAfterExpiry
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

