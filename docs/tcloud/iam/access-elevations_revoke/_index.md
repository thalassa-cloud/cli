---
linkTitle: "tcloud iam access-elevations revoke"
title: "iam access-elevations revoke"
slug: tcloud_iam_access-elevations_revoke
url: /docs/tcloud/iam/access-elevations_revoke/
weight: 9912
cascade:
  type: docs
---
## tcloud iam access-elevations revoke

Revoke an approved access elevation

```
tcloud iam access-elevations revoke <identity> [flags]
```

### Options

```
      --force                Skip the confirmation prompt
  -h, --help                 help for revoke
      --review-note string   Optional note attached to the review
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

