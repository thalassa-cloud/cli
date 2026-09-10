---
linkTitle: "tcloud quotas request-increase"
title: "quotas request-increase"
slug: tcloud_quotas_request-increase
url: /docs/tcloud/quotas/request-increase/
weight: 9726
cascade:
  type: docs
---
## tcloud quotas request-increase

Request a quota limit increase

```
tcloud quotas request-increase <name> [flags]
```

### Options

```
  -h, --help                help for request-increase
      --new-max-usage int   Requested maximum usage limit
      --reason string       Reason for the increase request
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

