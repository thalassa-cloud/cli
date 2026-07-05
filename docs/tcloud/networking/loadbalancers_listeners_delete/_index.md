---
linkTitle: "tcloud networking loadbalancers listeners delete"
title: "networking loadbalancers listeners delete"
slug: tcloud_networking_loadbalancers_listeners_delete
url: /docs/tcloud/networking/loadbalancers_listeners_delete/
weight: 9867
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners delete

Delete listener(s)

### Synopsis

Delete one or more listeners from a load balancer.

```
tcloud networking loadbalancers listeners delete LISTENER [LISTENER...] [flags]
```

### Examples

```
tcloud networking loadbalancers listeners delete listener-123 --loadbalancer lb-123 --force
```

### Options

```
      --force                 Skip confirmation
  -h, --help                  help for delete
      --loadbalancer string   Load balancer identity
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

