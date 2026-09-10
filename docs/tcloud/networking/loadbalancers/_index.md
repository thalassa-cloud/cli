---
linkTitle: "tcloud networking loadbalancers"
title: "networking loadbalancers"
slug: tcloud_networking_loadbalancers
url: /docs/tcloud/networking/loadbalancers/
weight: 9803
cascade:
  type: docs
---
## tcloud networking loadbalancers

Manage load balancers

### Synopsis

Manage load balancers, listeners, and related networking resources within the Thalassa Cloud Platform.

### Examples

```
tcloud networking loadbalancers list
tcloud networking loadbalancers create --name web --subnet subnet-123
tcloud networking loadbalancers listeners list --loadbalancer lb-123
```

### Options

```
  -h, --help   help for loadbalancers
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

* [tcloud networking](/docs/tcloud/tcloud_networking/)	 - Manage networking resources
* [tcloud networking loadbalancers create](/docs/tcloud/networking/loadbalancers_create/)	 - Create a load balancer
* [tcloud networking loadbalancers delete](/docs/tcloud/networking/loadbalancers_delete/)	 - Delete load balancer(s)
* [tcloud networking loadbalancers list](/docs/tcloud/networking/loadbalancers_list/)	 - List load balancers
* [tcloud networking loadbalancers listeners](/docs/tcloud/networking/loadbalancers_listeners/)	 - Manage load balancer listeners
* [tcloud networking loadbalancers update](/docs/tcloud/networking/loadbalancers_update/)	 - Update a load balancer
* [tcloud networking loadbalancers view](/docs/tcloud/networking/loadbalancers_view/)	 - View load balancer details

