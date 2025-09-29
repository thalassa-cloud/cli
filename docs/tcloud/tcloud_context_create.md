---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud context create"
title: "context create"
slug: tcloud_context_create
url: /docs/tcloud/tcloud_context_create/
weight: 9991
---
## tcloud context create

Create a new context with authentication and organisation

```
tcloud context create [flags]
```

### Options

```
      --create-context   creates a context (default true)
  -h, --help             help for create
      --name string      name of the context (default "default")
```

### Options inherited from parent commands

```
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

