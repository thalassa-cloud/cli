---
linkTitle: "tcloud networking routetables"
title: "networking routetables"
slug: tcloud_networking_routetables
url: /docs/tcloud/networking/routetables/
weight: 9776
cascade:
  type: docs
---
## tcloud networking routetables

Manage route tables

### Synopsis

Manage VPC route tables and their routes within the Thalassa Cloud Platform.

### Examples

```
tcloud networking routetables list
tcloud networking routetables create --name custom --vpc vpc-123
tcloud networking routetables routes list rt-123
```

### Options

```
  -h, --help   help for routetables
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
* [tcloud networking routetables create](/docs/tcloud/networking/routetables_create/)	 - Create a route table
* [tcloud networking routetables delete](/docs/tcloud/networking/routetables_delete/)	 - Delete route table(s)
* [tcloud networking routetables list](/docs/tcloud/networking/routetables_list/)	 - Get a list of routetables
* [tcloud networking routetables routes](/docs/tcloud/networking/routetables_routes/)	 - Manage routes in a route table
* [tcloud networking routetables update](/docs/tcloud/networking/routetables_update/)	 - Update a route table
* [tcloud networking routetables view](/docs/tcloud/networking/routetables_view/)	 - View a route table and its routes

