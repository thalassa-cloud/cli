---
linkTitle: "tcloud networking loadbalancers listeners list"
title: "networking loadbalancers listeners list"
slug: tcloud_networking_loadbalancers_listeners_list
url: /docs/tcloud/networking/loadbalancers_listeners_list/
weight: 9809
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners list

List load balancer listeners

### Synopsis

List listeners for a load balancer.

```
tcloud networking loadbalancers listeners list [flags]
```

### Examples

```
tcloud networking loadbalancers listeners list --loadbalancer lb-123
```

### Options

```
      --exact-time            Show exact time instead of relative time
  -h, --help                  help for list
      --loadbalancer string   Load balancer identity
      --no-header             Do not print the header
  -l, --selector string       Label selector (format: key1=value1,key2=value2)
      --show-labels           Show labels
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

