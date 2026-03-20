---
linkTitle: "tcloud dbaas backup list"
title: "dbaas backup list"
slug: tcloud_dbaas_backup_list
url: /docs/tcloud/dbaas/backup_list/
weight: 9973
cascade:
  type: docs
---
## tcloud dbaas backup list

List database backups

### Synopsis

List database backups for a specific cluster or all backups in the organisation

```
tcloud dbaas backup list [flags]
```

### Options

```
      --cluster string      Filter by database cluster identity, slug, or name
      --exact-time          Show exact time instead of relative time
  -h, --help                help for list
      --newer-than string   Filter backups newer than the specified duration (e.g., 7d, 1w, 1mo, 1y, 24h)
      --no-header           Do not print the header
      --older-than string   Filter backups older than the specified duration (e.g., 30d, 1w, 1mo, 1y, 24h)
  -l, --selector string     Label selector to filter backups (format: key1=value1,key2=value2)
      --show-labels         Show labels
      --status strings      Filter by backup status (can be specified multiple times, e.g., --status ready --status failed)
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

