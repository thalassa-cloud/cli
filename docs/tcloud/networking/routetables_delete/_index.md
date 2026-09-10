---
linkTitle: "tcloud networking routetables delete"
title: "networking routetables delete"
slug: tcloud_networking_routetables_delete
url: /docs/tcloud/networking/routetables_delete/
weight: 9787
cascade:
  type: docs
---
## tcloud networking routetables delete

Delete route table(s)

```
tcloud networking routetables delete ROUTE_TABLE [ROUTE_TABLE...] [flags]
```

### Options

```
      --force   Skip confirmation
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

* [tcloud networking routetables](/docs/tcloud/networking/routetables/)	 - Manage route tables

