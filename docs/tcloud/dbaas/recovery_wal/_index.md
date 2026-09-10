---
linkTitle: "tcloud dbaas recovery wal"
title: "dbaas recovery wal"
slug: tcloud_dbaas_recovery_wal
url: /docs/tcloud/dbaas/recovery_wal/
weight: 9947
cascade:
  type: docs
---
## tcloud dbaas recovery wal

Inspect WAL archives for recovery

### Synopsis

Browse WAL hierarchy, list WAL segments, and search WAL objects in a backup store

### Options

```
  -h, --help   help for wal
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
* [tcloud dbaas recovery wal hierarchy](/docs/tcloud/dbaas/recovery_wal_hierarchy/)	 - Show WAL archive hierarchy
* [tcloud dbaas recovery wal list](/docs/tcloud/dbaas/recovery_wal_list/)	 - List WAL segment objects
* [tcloud dbaas recovery wal search](/docs/tcloud/dbaas/recovery_wal_search/)	 - Search WAL archive objects

