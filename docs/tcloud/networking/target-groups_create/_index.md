---
linkTitle: "tcloud networking target-groups create"
title: "networking target-groups create"
slug: tcloud_networking_target-groups_create
url: /docs/tcloud/networking/target-groups_create/
weight: 9761
cascade:
  type: docs
---
## tcloud networking target-groups create

Create a target group

### Synopsis

Create a new load balancer target group.

```
tcloud networking target-groups create [flags]
```

### Examples

```
tcloud networking target-groups create --name web --vpc vpc-123 --port 8080 --protocol http
```

### Options

```
      --annotations strings           Annotations in key=value format
      --description string            Description of the target group
      --enable-proxy-protocol         Enable proxy protocol
  -h, --help                          help for create
      --labels strings                Labels in key=value format
      --loadbalancing-policy string   Load balancing policy (ROUND_ROBIN, RANDOM, MAGLEV)
      --name string                   Name of the target group
      --port int                      Target port
      --protocol string               Target protocol (tcp, udp, http, https, grpc, quic)
      --target-selector strings       Label selector for automatic target assignment (key=value)
      --vpc string                    VPC identity, slug, or name
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

* [tcloud networking target-groups](/docs/tcloud/networking/target-groups/)	 - Manage load balancer target groups

