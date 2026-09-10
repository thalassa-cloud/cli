---
linkTitle: "tcloud dbaas recovery wal hierarchy"
title: "dbaas recovery wal hierarchy"
slug: tcloud_dbaas_recovery_wal_hierarchy
url: /docs/tcloud/dbaas/recovery_wal_hierarchy/
weight: 9950
cascade:
  type: docs
---
## tcloud dbaas recovery wal hierarchy

Show WAL archive hierarchy

### Synopsis

Show WAL timelines and log folders for a cluster in a backup store

```
tcloud dbaas recovery wal hierarchy [cluster] [flags]
```

### Options

```
      --backup-store string   Backup store identity (optional when the cluster has a backup store)
      --cluster string        Cluster identity (alternative to the positional argument)
      --exact-time            Show exact time instead of relative time
  -h, --help                  help for hierarchy
      --kinds strings         Filter by WAL object kinds (wal, backup_history, timeline_history)
      --no-header             Do not print the header
      --timeline string       Filter to a single timeline
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

