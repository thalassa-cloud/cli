---
linkTitle: "tcloud dbaas backup delete"
title: "dbaas backup delete"
slug: tcloud_dbaas_backup_delete
url: /docs/tcloud/dbaas/backup_delete/
weight: 9974
cascade:
  type: docs
---
## tcloud dbaas backup delete

Delete database backup(s)

### Synopsis

Delete database backup(s) by identity, label selector, or all failed backups

```
tcloud dbaas backup delete [flags]
```

### Options

```
      --all-failed        Delete all failed backups
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter backups (format: key1=value1,key2=value2)
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

