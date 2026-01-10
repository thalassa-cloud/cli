---
linkTitle: "tcloud object-storage create"
title: "object-storage create"
slug: tcloud_object-storage_create
url: /docs/tcloud/object-storage/create/
weight: 9927
cascade:
  type: docs
---
## tcloud object-storage create

Create an object storage bucket

```
tcloud object-storage create [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format
  -h, --help                  help for create
      --labels strings        Labels in key=value format
      --name string           Name of the bucket
      --no-header             Do not print the header
      --object-lock-enabled   Enable object lock
      --policy string         Bucket policy as JSON string or path to a JSON file
      --region string         Region for the bucket
      --timeout string        Timeout for waiting (e.g., 5m, 10m, 1h, default: 10m)
      --versioning            Enable versioning for the bucket
  -w, --wait                  Wait for the bucket to be ready
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

