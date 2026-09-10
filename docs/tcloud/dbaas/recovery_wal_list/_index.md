---
linkTitle: "tcloud dbaas recovery wal list"
title: "dbaas recovery wal list"
slug: tcloud_dbaas_recovery_wal_list
url: /docs/tcloud/dbaas/recovery_wal_list/
weight: 9949
cascade:
  type: docs
---
## tcloud dbaas recovery wal list

List WAL segment objects

### Synopsis

List WAL segment objects for a timeline/log folder, or as a flat listing with --flat

```
tcloud dbaas recovery wal list [cluster] [flags]
```

### Options

```
      --backup-store string         Backup store identity (optional when the cluster has a backup store)
      --cluster string              Cluster identity (alternative to the positional argument)
      --continuation-token string   Continuation token from a previous listing
      --exact-time                  Show exact time instead of relative time
      --flat                        List WAL objects without requiring timeline and log
  -h, --help                        help for list
      --kinds strings               Filter by WAL object kinds (wal, backup_history, timeline_history)
      --limit int                   Maximum number of objects to return
      --log string                  Log hex (required unless --flat)
      --no-header                   Do not print the header
      --timeline string             Timeline hex (required unless --flat)
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

