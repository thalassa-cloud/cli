---
linkTitle: "tcloud object-storage delete"
title: "object-storage delete"
slug: tcloud_object-storage_delete
url: /docs/tcloud/object-storage/delete/
weight: 9926
cascade:
  type: docs
---
## tcloud object-storage delete

Delete an object storage bucket

### Synopsis

Delete an object storage bucket by name. This will permanently delete the bucket and all its contents.

```
tcloud object-storage delete [flags]
```

### Examples

```
tcloud storage object-storage delete my-bucket
tcloud storage object-storage delete my-bucket --force
```

### Options

```
      --force            Force the deletion and skip the confirmation
  -h, --help             help for delete
      --timeout string   Timeout for waiting (e.g., 5m, 10m, 1h, default: 10m)
  -w, --wait             Wait for the bucket to be deleted
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

* [tcloud object-storage](/docs/tcloud/tcloud_object-storage/)	 - Manage object storage buckets

