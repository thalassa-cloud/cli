---
linkTitle: "tcloud networking loadbalancers update"
title: "networking loadbalancers update"
slug: tcloud_networking_loadbalancers_update
url: /docs/tcloud/networking/loadbalancers_update/
weight: 9805
cascade:
  type: docs
---
## tcloud networking loadbalancers update

Update a load balancer

### Synopsis

Update properties of an existing load balancer.

```
tcloud networking loadbalancers update LOADBALANCER [flags]
```

### Examples

```
tcloud networking loadbalancers update lb-123 --name web-prod
tcloud networking loadbalancers update lb-123 --delete-protection
```

### Options

```
      --annotations strings       Annotations in key=value format
      --delete-protection         Enable delete protection
      --description string        Description of the load balancer
  -h, --help                      help for update
      --labels strings            Labels in key=value format
      --name string               Name of the load balancer
      --security-groups strings   Security group identities to attach
      --subnet string             Subnet identity, slug, or name
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

