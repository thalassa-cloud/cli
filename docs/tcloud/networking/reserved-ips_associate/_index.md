---
linkTitle: "tcloud networking reserved-ips associate"
title: "networking reserved-ips associate"
slug: tcloud_networking_reserved-ips_associate
url: /docs/tcloud/networking/reserved-ips_associate/
weight: 9796
cascade:
  type: docs
---
## tcloud networking reserved-ips associate

Associate a reserved IP with a resource

### Synopsis

Associate a reserved IP with exactly one of a load balancer or NAT gateway.

```
tcloud networking reserved-ips associate [flags]
```

### Examples

```
tcloud networking reserved-ips associate rip-123 --loadbalancer lb-456
tcloud networking reserved-ips associate rip-123 --nat-gateway ngw-789
```

### Options

```
  -h, --help                  help for associate
      --loadbalancer string   Load balancer identity to associate with
      --nat-gateway string    NAT gateway identity to associate with
      --no-header             Do not print the header
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

