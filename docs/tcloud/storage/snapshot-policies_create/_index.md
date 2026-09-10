---
linkTitle: "tcloud storage snapshot-policies create"
title: "storage snapshot-policies create"
slug: tcloud_storage_snapshot-policies_create
url: /docs/tcloud/storage/snapshot-policies_create/
weight: 9693
cascade:
  type: docs
---
## tcloud storage snapshot-policies create

Create a snapshot policy

### Synopsis

Create a new automated snapshot policy for volumes in a region.

```
tcloud storage snapshot-policies create [flags]
```

### Examples

```
tcloud storage snapshot-policies create --name daily --region nl-ams --schedule "0 2 * * *" --ttl 168h --timezone UTC --target-type selector --selector backup=true
tcloud storage snapshot-policies create --name weekly --region nl-ams --schedule "0 3 * * 0" --ttl 720h --target-type explicit --volumes vol-1,vol-2
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description of the snapshot policy
      --enabled               Enable the snapshot policy (default true)
  -h, --help                  help for create
      --keep-count int        Maximum number of snapshots to retain
      --labels strings        Labels in key=value format
      --name string           Name of the snapshot policy
      --no-header             Do not print the header
      --region string         Region of the snapshot policy
      --schedule string       Cron schedule for snapshot creation
      --selector strings      Label selectors for target volumes (key=value)
      --target-type string    Target type: selector or explicit
      --timezone string       Timezone for the schedule (default "UTC")
      --ttl string            Snapshot retention duration (e.g. 24h, 168h)
      --volumes strings       Volume identities when target-type is explicit
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

* [tcloud storage snapshot-policies](/docs/tcloud/storage/snapshot-policies/)	 - Manage snapshot policies

