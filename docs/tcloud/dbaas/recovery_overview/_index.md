---
linkTitle: "tcloud dbaas recovery overview"
title: "dbaas recovery overview"
slug: tcloud_dbaas_recovery_overview
url: /docs/tcloud/dbaas/recovery_overview/
weight: 9951
cascade:
  type: docs
---
## tcloud dbaas recovery overview

Show backup recovery and PITR overview

### Synopsis

Show aggregated recovery information for a database cluster or backup store.

The overview includes backup counts, the approximate PITR window, WAL chain coverage
per backup, and any gaps detected in the archive.

```
tcloud dbaas recovery overview [cluster] [flags]
```

### Examples

```
tcloud dbaas recovery overview dbc-123
tcloud dbaas recovery overview --backup-store bst-123
```

### Options

```
      --backup string         Limit the overview to a single backup identity
      --backup-store string   Backup store identity (optional when a cluster is provided)
      --exact-time            Show exact time instead of relative time
      --exclude-unavailable   Exclude unavailable backups from the overview
  -h, --help                  help for overview
      --include-unavailable   Include unavailable backups in the overview (default true)
      --no-header             Do not print the header
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

* [tcloud dbaas recovery](/docs/tcloud/dbaas/recovery/)	 - Inspect database backup recovery and PITR windows

