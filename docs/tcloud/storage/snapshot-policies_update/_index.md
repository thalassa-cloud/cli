---
linkTitle: "tcloud storage snapshot-policies update"
title: "storage snapshot-policies update"
slug: tcloud_storage_snapshot-policies_update
url: /docs/tcloud/storage/snapshot-policies_update/
weight: 9690
cascade:
  type: docs
---
## tcloud storage snapshot-policies update

Update a snapshot policy

### Synopsis

Update an existing snapshot policy. Unspecified fields are preserved from the current resource.

```
tcloud storage snapshot-policies update [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description of the snapshot policy
      --enabled               Enable or disable the snapshot policy (default true)
  -h, --help                  help for update
      --keep-count int        Maximum number of snapshots to retain
      --labels strings        Labels in key=value format
      --name string           Name of the snapshot policy
      --no-header             Do not print the header
      --schedule string       Cron schedule for snapshot creation
      --selector strings      Label selectors for target volumes (key=value)
      --target-type string    Target type: selector or explicit
      --timezone string       Timezone for the schedule
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

