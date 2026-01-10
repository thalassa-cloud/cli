---
linkTitle: "tcloud storage volumes create"
title: "storage volumes create"
slug: tcloud_storage_volumes_create
url: /docs/tcloud/storage/volumes_create/
weight: 9913
cascade:
  type: docs
---
## tcloud storage volumes create

Create a volume

### Synopsis

Create a new storage volume. The volume can be attached to machines after creation.

```
tcloud storage volumes create [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format (can be specified multiple times)
      --delete-protection     Enable delete protection
      --description string    Description of the volume
  -h, --help                  help for create
      --labels strings        Labels in key=value format (can be specified multiple times)
      --name string           Name of the volume (required)
      --no-header             Do not print the header
      --region string         Region of the volume (required)
      --size int              Size of the volume in GB (required)
      --type string           Volume type (default "block")
      --wait                  Wait for the volume to be available before returning
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

* [tcloud storage volumes](/docs/tcloud/storage/volumes/)	 - Manage storage volumes

