---
linkTitle: "tcloud dbaas backup-schedules list"
title: "dbaas backup-schedules list"
slug: tcloud_dbaas_backup-schedules_list
url: /docs/tcloud/dbaas/backup-schedules_list/
weight: 9968
cascade:
  type: docs
---
## tcloud dbaas backup-schedules list

List database backup schedules

### Synopsis

List database backup schedules for a specific cluster or all schedules in the organisation

```
tcloud dbaas backup-schedules list [flags]
```

### Options

```
      --cluster string   Filter by database cluster identity, slug, or name
      --exact-time       Show exact time instead of relative time
  -h, --help             help for list
      --no-header        Do not print the header
      --show-labels      Show labels
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

* [tcloud dbaas backup-schedules](/docs/tcloud/dbaas/backup-schedules/)	 - Manage database backup schedules

