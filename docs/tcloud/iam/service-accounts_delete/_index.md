---
linkTitle: "tcloud iam service-accounts delete"
title: "iam service-accounts delete"
slug: tcloud_iam_service-accounts_delete
url: /docs/tcloud/iam/service-accounts_delete/
weight: 9923
cascade:
  type: docs
---
## tcloud iam service-accounts delete

Delete a service account

```
tcloud iam service-accounts delete <identity> [flags]
```

### Options

```
      --force       Skip the confirmation prompt and delete
  -h, --help        help for delete
      --no-header   Do not print table headers
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

* [tcloud iam service-accounts](/docs/tcloud/iam/service-accounts/)	 - Organisation service accounts

