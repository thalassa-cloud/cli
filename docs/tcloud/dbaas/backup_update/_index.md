---
linkTitle: "tcloud dbaas backup update"
title: "dbaas backup update"
slug: tcloud_dbaas_backup_update
url: /docs/tcloud/dbaas/backup_update/
weight: 9970
cascade:
  type: docs
---
## tcloud dbaas backup update

Update a database backup

### Synopsis

Update backup properties such as delete protection

```
tcloud dbaas backup update [flags]
```

### Options

```
      --delete-protection   Enable or disable delete protection
  -h, --help                help for update
      --no-header           Do not print the header
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

