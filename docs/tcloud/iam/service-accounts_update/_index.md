---
linkTitle: "tcloud iam service-accounts update"
title: "iam service-accounts update"
slug: tcloud_iam_service-accounts_update
url: /docs/tcloud/iam/service-accounts_update/
weight: 9922
cascade:
  type: docs
---
## tcloud iam service-accounts update

Update a service account (only set flags are sent)

```
tcloud iam service-accounts update <identity> [flags]
```

### Options

```
      --annotations strings   Replace annotations (key=value, repeatable)
      --description string    Description (empty to clear)
  -h, --help                  help for update
      --labels strings        Replace labels (key=value, repeatable)
      --name string           Name
      --no-header             Do not print table headers
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

* [tcloud iam service-accounts](/docs/tcloud/iam/service-accounts/)	 - Organisation service accounts

