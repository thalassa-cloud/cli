---
linkTitle: "tcloud networking routetables create"
title: "networking routetables create"
slug: tcloud_networking_routetables_create
url: /docs/tcloud/networking/routetables_create/
weight: 9788
cascade:
  type: docs
---
## tcloud networking routetables create

Create a route table

```
tcloud networking routetables create [flags]
```

### Examples

```
tcloud networking routetables create --name private --vpc vpc-123
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Description
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --name string           Name of the route table
      --vpc string            VPC identity
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

* [tcloud networking routetables](/docs/tcloud/networking/routetables/)	 - Manage route tables

