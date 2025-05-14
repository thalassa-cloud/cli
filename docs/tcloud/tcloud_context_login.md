---
date: 2025-05-14T17:58:13+02:00
linkTitle: "tcloud context login"
title: "context login"
slug: tcloud_context_login
url: /docs/tcloud/tcloud_context_login/
weight: 9987
---
## tcloud context login

Login to Thalassa Cloud

### Synopsis

Login to Thalassa Cloud using a personal access token, using the current context. Overrides the current context if --name is set.

```
tcloud context login [flags]
```

### Options

```
  -h, --help          help for login
      --name string   name of the context (default "default")
```

### Options inherited from parent commands

```
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

