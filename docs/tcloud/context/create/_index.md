---
linkTitle: "tcloud context create"
title: "context create"
slug: tcloud_context_create
url: /docs/tcloud/context/create/
weight: 9987
cascade:
  type: docs
---
## tcloud context create

Create a new context with authentication and organisation

```
tcloud context create [flags]
```

### Options

```
      --browser          log in through the browser using OIDC (default when no other credentials are provided)
      --create-context   creates a context (default true)
  -h, --help             help for create
      --name string      name of the context (default "default")
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

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

