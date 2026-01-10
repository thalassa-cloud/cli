---
linkTitle: "tcloud object-storage update"
title: "object-storage update"
slug: tcloud_object-storage_update
url: /docs/tcloud/object-storage/update/
weight: 9924
cascade:
  type: docs
---
## tcloud object-storage update

Update an object storage bucket

```
tcloud object-storage update [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format
  -h, --help                  help for update
      --labels strings        Labels in key=value format
      --no-header             Do not print the header
      --object-lock-enabled   Enable object lock
      --policy string         Bucket policy as JSON string or path to a JSON file
      --public                Make the bucket publicly accessible
      --versioning            Enable versioning for the bucket
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

