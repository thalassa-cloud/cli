---
linkTitle: "tcloud quotas list"
title: "quotas list"
slug: tcloud_quotas_list
url: /docs/tcloud/quotas/list/
weight: 9813
cascade:
  type: docs
---
## tcloud quotas list

List organisation quotas

```
tcloud quotas list [flags]
```

### Options

```
  -h, --help                     help for list
      --no-header                Do not print table headers
      --show-increase-requests   Include requested increase limits and their decision status
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

* [tcloud quotas](/docs/tcloud/tcloud_quotas/)	 - View and request changes to organisation resource quotas

