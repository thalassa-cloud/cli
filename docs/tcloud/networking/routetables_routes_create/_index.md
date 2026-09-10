---
linkTitle: "tcloud networking routetables routes create"
title: "networking routetables routes create"
slug: tcloud_networking_routetables_routes_create
url: /docs/tcloud/networking/routetables_routes_create/
weight: 9785
cascade:
  type: docs
---
## tcloud networking routetables routes create

Create a route in a route table

```
tcloud networking routetables routes create ROUTE_TABLE [flags]
```

### Options

```
      --destination string       Destination CIDR block
      --gateway string           Target gateway identity
      --gateway-address string   Gateway address
  -h, --help                     help for create
      --nat-gateway string       Target NAT gateway identity
      --no-header                Do not print the header
      --vpc-peering string       Target VPC peering connection identity
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

* [tcloud networking routetables routes](/docs/tcloud/networking/routetables_routes/)	 - Manage routes in a route table

