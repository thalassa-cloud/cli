---
linkTitle: "tcloud iam access-elevations get"
title: "iam access-elevations get"
slug: tcloud_iam_access-elevations_get
url: /docs/tcloud/iam/access-elevations_get/
weight: 9915
cascade:
  type: docs
---
## tcloud iam access-elevations get

Show an access elevation request

```
tcloud iam access-elevations get <identity> [flags]
```

### Options

```
      --exact-time   Show full timestamps instead of relative time
  -h, --help         help for get
      --mine         Fetch via your own requests API
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

