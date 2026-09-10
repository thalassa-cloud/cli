---
linkTitle: "tcloud dbaas recovery wal search"
title: "dbaas recovery wal search"
slug: tcloud_dbaas_recovery_wal_search
url: /docs/tcloud/dbaas/recovery_wal_search/
weight: 9948
cascade:
  type: docs
---
## tcloud dbaas recovery wal search

Search WAL archive objects

### Synopsis

Search WAL archive objects by WAL name or hex substring

```
tcloud dbaas recovery wal search [cluster] [flags]
```

### Options

```
      --backup-store string   Backup store identity (optional when the cluster has a backup store)
      --cluster string        Cluster identity (alternative to the positional argument)
      --exact-time            Show exact time instead of relative time
  -h, --help                  help for search
      --kinds strings         Filter by WAL object kinds (wal, backup_history, timeline_history)
      --limit int             Maximum number of objects to return
      --no-header             Do not print the header
      --query string          WAL name or hex substring to search for (required)
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

* [tcloud dbaas recovery wal](/docs/tcloud/dbaas/recovery_wal/)	 - Inspect WAL archives for recovery

