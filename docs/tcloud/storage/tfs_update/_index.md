---
linkTitle: "tcloud storage tfs update"
title: "storage tfs update"
slug: tcloud_storage_tfs_update
url: /docs/tcloud/storage/tfs_update/
weight: 9841
cascade:
  type: docs
---
## tcloud storage tfs update

Update a TFS instance

### Synopsis

Update properties of an existing TFS instance.

```
tcloud storage tfs update [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format (can be specified multiple times)
      --delete-protection     Enable or disable delete protection
      --description string    Description of the TFS instance
  -h, --help                  help for update
      --labels strings        Labels in key=value format (can be specified multiple times)
      --name string           Name of the TFS instance
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud storage tfs](/docs/tcloud/storage/tfs/)	 - Manage TFS (Thalassa File System) instances

