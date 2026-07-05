---
linkTitle: "tcloud networking loadbalancers listeners update"
title: "networking loadbalancers listeners update"
slug: tcloud_networking_loadbalancers_listeners_update
url: /docs/tcloud/networking/loadbalancers_listeners_update/
weight: 9865
cascade:
  type: docs
---
## tcloud networking loadbalancers listeners update

Update a listener

### Synopsis

Update properties of an existing load balancer listener.

```
tcloud networking loadbalancers listeners update LISTENER [flags]
```

### Examples

```
tcloud networking loadbalancers listeners update listener-123 --loadbalancer lb-123 --port 8080
```

### Options

```
      --allowed-sources strings          Allowed source CIDR blocks
      --annotations strings              Annotations in key=value format
      --connection-idle-timeout uint32   Connection idle timeout in seconds
      --description string               Description of the listener
  -h, --help                             help for update
      --labels strings                   Labels in key=value format
      --loadbalancer string              Load balancer identity
      --max-connections uint32           Maximum connections
      --name string                      Name of the listener
      --port int                         Listener port
      --protocol string                  Listener protocol
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

