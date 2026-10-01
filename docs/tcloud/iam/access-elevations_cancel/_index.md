---
linkTitle: "tcloud iam access-elevations cancel"
title: "iam access-elevations cancel"
slug: tcloud_iam_access-elevations_cancel
url: /docs/tcloud/iam/access-elevations_cancel/
weight: 9917
cascade:
  type: docs
---
## tcloud iam access-elevations cancel

Cancel your pending access elevation request

```
tcloud iam access-elevations cancel <identity> [flags]
```

### Options

```
      --force   Skip the confirmation prompt
  -h, --help    help for cancel
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam access-elevations](/docs/tcloud/iam/access-elevations/)	 - Request and review temporary IAM access elevations

