---
linkTitle: "tcloud networking routetables routes set"
title: "networking routetables routes set"
slug: tcloud_networking_routetables_routes_set
url: /docs/tcloud/networking/routetables_routes_set/
weight: 9782
cascade:
  type: docs
---
## tcloud networking routetables routes set

Replace all routes in a route table from a JSON file

### Synopsis

Batch-update routes for a route table. The file must contain a JSON array of UpdateRouteTableRoute objects.

```
tcloud networking routetables routes set ROUTE_TABLE [flags]
```

### Examples

```
tcloud networking routetables routes set rt-123 --file routes.json
```

### Options

```
      --file string   JSON file containing an array of routes
  -h, --help          help for set
      --no-header     Do not print the header
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

