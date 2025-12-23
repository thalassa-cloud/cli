---
linkTitle: "tcloud storage snapshots create"
title: "storage snapshots create"
slug: tcloud_storage_snapshots_create
url: /docs/tcloud/storage/snapshots_create/
weight: 9940
cascade:
  type: docs
---
## tcloud storage snapshots create

Create a snapshot

### Synopsis

Create a snapshot from a volume by providing a name and volume ID.

```
tcloud storage snapshots create <name> [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format (can be specified multiple times)
      --delete-protection     Enable delete protection for the snapshot
      --description string    Description of the snapshot
  -h, --help                  help for create
      --labels strings        Labels in key=value format (can be specified multiple times)
      --volume string         Volume identity to create snapshot from
      --wait                  Wait for the snapshot to be ready for use
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

* [tcloud storage snapshots](/docs/tcloud/storage/snapshots/)	 - Manage volume snapshots

