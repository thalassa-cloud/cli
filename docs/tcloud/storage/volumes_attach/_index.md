---
linkTitle: "tcloud storage volumes attach"
title: "storage volumes attach"
slug: tcloud_storage_volumes_attach
url: /docs/tcloud/storage/volumes_attach/
weight: 9914
cascade:
  type: docs
---
## tcloud storage volumes attach

Attach volume(s) to a virtual machine

### Synopsis

Attach one or more volumes to a virtual machine instance by identity.

```
tcloud storage volumes attach <volume-id> [<volume-id> ...] [flags]
```

### Options

```
  -h, --help              help for attach
      --instance string   Virtual machine instance identity
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

