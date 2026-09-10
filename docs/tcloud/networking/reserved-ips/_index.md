---
linkTitle: "tcloud networking reserved-ips"
title: "networking reserved-ips"
slug: tcloud_networking_reserved-ips
url: /docs/tcloud/networking/reserved-ips/
weight: 9789
cascade:
  type: docs
---
## tcloud networking reserved-ips

Manage reserved IP addresses

### Synopsis

Manage reserved public IP addresses that can be associated with load balancers or NAT gateways.

### Examples

```
tcloud networking reserved-ips list
tcloud networking reserved-ips create --name my-ip --region nl-ams
tcloud networking reserved-ips associate rip-123 --nat-gateway ngw-456
```

### Options

```
  -h, --help   help for reserved-ips
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
* [tcloud networking reserved-ips associate](/docs/tcloud/networking/reserved-ips_associate/)	 - Associate a reserved IP with a resource
* [tcloud networking reserved-ips create](/docs/tcloud/networking/reserved-ips_create/)	 - Create a reserved IP address
* [tcloud networking reserved-ips delete](/docs/tcloud/networking/reserved-ips_delete/)	 - Delete reserved IP address(es)
* [tcloud networking reserved-ips disassociate](/docs/tcloud/networking/reserved-ips_disassociate/)	 - Disassociate a reserved IP from its resource
* [tcloud networking reserved-ips list](/docs/tcloud/networking/reserved-ips_list/)	 - List reserved IP addresses
* [tcloud networking reserved-ips update](/docs/tcloud/networking/reserved-ips_update/)	 - Update a reserved IP address
* [tcloud networking reserved-ips view](/docs/tcloud/networking/reserved-ips_view/)	 - View reserved IP details

