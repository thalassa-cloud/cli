---
linkTitle: "tcloud networking routetables routes delete"
title: "networking routetables routes delete"
slug: tcloud_networking_routetables_routes_delete
url: /docs/tcloud/networking/routetables_routes_delete/
weight: 9784
cascade:
  type: docs
---
## tcloud networking routetables routes delete

Delete a route from a route table

```
tcloud networking routetables routes delete ROUTE_TABLE ROUTE [flags]
```

### Options

```
      --force   Force the deletion and skip the confirmation
  -h, --help    help for delete
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

