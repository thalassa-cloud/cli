---
linkTitle: "tcloud networking loadbalancers listeners create"
title: "networking loadbalancers listeners create"
slug: tcloud_networking_loadbalancers_listeners_create
url: /docs/tcloud/networking/loadbalancers_listeners_create/
weight: 9811
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners create

Create a listener

### Synopsis

Create a listener on a load balancer.

```
tcloud networking loadbalancers listeners create [flags]
```

### Examples

```
tcloud networking loadbalancers listeners create --loadbalancer lb-123 --name http --port 80 --protocol tcp --target-group tg-123
```

### Options

```
      --allowed-sources strings          Allowed source CIDR blocks
      --annotations strings              Annotations in key=value format
      --connection-idle-timeout uint32   Connection idle timeout in seconds
      --description string               Description of the listener
  -h, --help                             help for create
      --labels strings                   Labels in key=value format
      --loadbalancer string              Load balancer identity
      --max-connections uint32           Maximum connections
      --name string                      Name of the listener
      --port int                         Listener port
      --protocol string                  Listener protocol (tcp, udp)
      --target-group string              Target group identity
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

