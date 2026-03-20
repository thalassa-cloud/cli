---
linkTitle: "tcloud storage tfs create"
title: "storage tfs create"
slug: tcloud_storage_tfs_create
url: /docs/tcloud/storage/tfs_create/
weight: 9844
cascade:
  type: docs
---
## tcloud storage tfs create

Create a TFS instance

### Synopsis

Create a new TFS (Thalassa File System) instance for shared file storage.

```
tcloud storage tfs create [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format (can be specified multiple times)
      --delete-protection     Enable delete protection
      --description string    Description of the TFS instance
  -h, --help                  help for create
      --labels strings        Labels in key=value format (can be specified multiple times)
      --name string           Name of the TFS instance (required)
      --no-header             Do not print the header
      --region string         Region of the TFS instance (required)
      --size int              Size of the TFS instance in GB (required) (default 1)
      --subnet string         Subnet of the TFS instance (required)
      --vpc string            VPC of the TFS instance (required)
      --wait                  Wait for the TFS instance to be available before returning
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

