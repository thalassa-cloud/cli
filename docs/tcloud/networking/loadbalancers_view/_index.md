---
linkTitle: "tcloud networking loadbalancers view"
title: "networking loadbalancers view"
slug: tcloud_networking_loadbalancers_view
url: /docs/tcloud/networking/loadbalancers_view/
weight: 9861
cascade:
  type: docs
---
## tcloud networking loadbalancers view

View load balancer details

### Synopsis

View detailed information about a load balancer.

```
tcloud networking loadbalancers view LOADBALANCER [flags]
```

### Examples

```
tcloud networking loadbalancers view lb-123
tcloud networking loadbalancers view lb-123 --output yaml
```

### Options

```
  -h, --help            help for view
  -o, --output string   Output format (yaml)
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

* [tcloud networking loadbalancers](/docs/tcloud/networking/loadbalancers/)	 - Manage load balancers

