---
linkTitle: "tcloud iam access-elevations list"
title: "iam access-elevations list"
slug: tcloud_iam_access-elevations_list
url: /docs/tcloud/iam/access-elevations_list/
weight: 9914
cascade:
  type: docs
---
## tcloud iam access-elevations list

List access elevation requests

### Synopsis

List access elevation requests visible to approvers in the current scope.

Use --mine to list only requests you created.

```
tcloud iam access-elevations list [flags]
```

### Options

```
  -h, --help        help for list
      --mine        List only requests you created
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
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam access-elevations](/docs/tcloud/iam/access-elevations/)	 - Request and review temporary IAM access elevations

