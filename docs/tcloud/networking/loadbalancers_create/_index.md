---
linkTitle: "tcloud networking loadbalancers create"
title: "networking loadbalancers create"
slug: tcloud_networking_loadbalancers_create
url: /docs/tcloud/networking/loadbalancers_create/
weight: 9871
cascade:
  type: docs
---
## tcloud networking loadbalancers create

Create a load balancer

### Synopsis

Create a new load balancer in the specified subnet.

```
tcloud networking loadbalancers create [flags]
```

### Examples

```
tcloud networking loadbalancers create --name web --subnet subnet-123
tcloud networking loadbalancers create --name internal --subnet subnet-123 --internal --wait
```

### Options

```
      --annotations strings       Annotations in key=value format
      --delete-protection         Enable delete protection
      --description string        Description of the load balancer
  -h, --help                      help for create
      --internal                  Create an internal load balancer (no public IP)
      --labels strings            Labels in key=value format
      --name string               Name of the load balancer
      --security-groups strings   Security group identities to attach
      --subnet string             Subnet identity, slug, or name
      --wait                      Wait for the load balancer to be ready
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

