---
linkTitle: "tcloud networking loadbalancers listeners view"
title: "networking loadbalancers listeners view"
slug: tcloud_networking_loadbalancers_listeners_view
url: /docs/tcloud/networking/loadbalancers_listeners_view/
weight: 9864
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners view

View listener details

### Synopsis

View detailed information about a load balancer listener.

```
tcloud networking loadbalancers listeners view LISTENER [flags]
```

### Examples

```
tcloud networking loadbalancers listeners view listener-123 --loadbalancer lb-123
```

### Options

```
  -h, --help                  help for view
      --loadbalancer string   Load balancer identity
  -o, --output string         Output format (yaml)
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

* [tcloud networking loadbalancers listeners](/docs/tcloud/networking/loadbalancers_listeners/)	 - Manage load balancer listeners

