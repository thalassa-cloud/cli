---
linkTitle: "tcloud networking loadbalancers listeners"
title: "networking loadbalancers listeners"
slug: tcloud_networking_loadbalancers_listeners
url: /docs/tcloud/networking/loadbalancers_listeners/
weight: 9806
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners

Manage load balancer listeners

### Synopsis

Manage listeners on a load balancer. All commands require --loadbalancer.

### Examples

```
tcloud networking loadbalancers listeners list --loadbalancer lb-123
tcloud networking loadbalancers listeners create --loadbalancer lb-123 --name http --port 80 --protocol http --target-group tg-123
```

### Options

```
  -h, --help   help for listeners
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
* [tcloud networking loadbalancers listeners create](/docs/tcloud/networking/loadbalancers_listeners_create/)	 - Create a listener
* [tcloud networking loadbalancers listeners delete](/docs/tcloud/networking/loadbalancers_listeners_delete/)	 - Delete listener(s)
* [tcloud networking loadbalancers listeners list](/docs/tcloud/networking/loadbalancers_listeners_list/)	 - List load balancer listeners
* [tcloud networking loadbalancers listeners update](/docs/tcloud/networking/loadbalancers_listeners_update/)	 - Update a listener
* [tcloud networking loadbalancers listeners view](/docs/tcloud/networking/loadbalancers_listeners_view/)	 - View listener details

