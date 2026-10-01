---
linkTitle: "tcloud iam policies get"
title: "iam policies get"
slug: tcloud_iam_policies_get
url: /docs/tcloud/iam/policies_get/
weight: 9893
cascade:
  type: docs
---
## tcloud iam policies get

Show a policy including rules and bindings summary

```
tcloud iam policies get <policy> [flags]
```

### Options

```
      --exact-time   Show full timestamps instead of relative time
  -h, --help         help for get
      --no-header    Do not print table headers
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

* [tcloud iam policies](/docs/tcloud/iam/policies/)	 - IAM policies, permission rules, and bindings

