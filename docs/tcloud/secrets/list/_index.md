---
linkTitle: "tcloud secrets list"
title: "secrets list"
slug: tcloud_secrets_list
url: /docs/tcloud/secrets/list/
weight: 9698
cascade:
  type: docs
---
## tcloud secrets list

List secrets under a path prefix (metadata only)

```
tcloud secrets list [flags]
```

### Options

```
      --exact-time      Show full timestamps instead of relative time
  -h, --help            help for list
      --no-header       Do not print table headers
      --prefix string   Path prefix (default "/")
      --region string   Region
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

* [tcloud secrets](/docs/tcloud/tcloud_secrets/)	 - Manage secrets (beta)

