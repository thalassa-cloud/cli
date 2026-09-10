---
linkTitle: "tcloud networking reserved-ips disassociate"
title: "networking reserved-ips disassociate"
slug: tcloud_networking_reserved-ips_disassociate
url: /docs/tcloud/networking/reserved-ips_disassociate/
weight: 9793
cascade:
  type: docs
---
## tcloud networking reserved-ips disassociate

Disassociate a reserved IP from its resource

### Synopsis

Detach a reserved IP from its currently associated load balancer or NAT gateway.

```
tcloud networking reserved-ips disassociate [flags]
```

### Examples

```
tcloud networking reserved-ips disassociate rip-123
```

### Options

```
  -h, --help        help for disassociate
      --no-header   Do not print the header
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

* [tcloud networking reserved-ips](/docs/tcloud/networking/reserved-ips/)	 - Manage reserved IP addresses

