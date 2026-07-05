---
linkTitle: "tcloud quotas"
title: "quotas"
slug: tcloud_quotas
url: /docs/tcloud/tcloud_quotas/
weight: 9811
cascade:
  type: docs
---
## tcloud quotas

View and request changes to organisation resource quotas

### Synopsis

List organisation quotas, inspect current usage, and submit increase requests.
Quotas apply per organisation in your current context.

### Examples

```
  # List all quotas for the current organisation
  tcloud quotas list

  # Show a single quota including pending increase requests
  tcloud quotas get machines

  # Request a higher limit
  tcloud quotas request-increase machines --new-max-usage 50 --reason "Growing production fleet"
```

### Options

```
  -h, --help   help for quotas
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud quotas get](/docs/tcloud/quotas/get/)	 - Show details for a quota
* [tcloud quotas list](/docs/tcloud/quotas/list/)	 - List organisation quotas
* [tcloud quotas request-increase](/docs/tcloud/quotas/request-increase/)	 - Request a quota limit increase

