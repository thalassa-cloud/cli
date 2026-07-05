---
linkTitle: "tcloud networking loadbalancers delete"
title: "networking loadbalancers delete"
slug: tcloud_networking_loadbalancers_delete
url: /docs/tcloud/networking/loadbalancers_delete/
weight: 9870
cascade:
  type: docs
---
## tcloud networking loadbalancers delete

Delete load balancer(s)

### Synopsis

Delete load balancer(s) by identity or label selector.

```
tcloud networking loadbalancers delete [LOADBALANCER...] [flags]
```

### Examples

```
tcloud networking loadbalancers delete lb-123
tcloud networking loadbalancers delete lb-123 --wait
tcloud networking loadbalancers delete --selector env=test --force
```

### Options

```
      --force             Force deletion and skip confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector (format: key1=value1,key2=value2)
      --wait              Wait for the load balancer(s) to be deleted
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

