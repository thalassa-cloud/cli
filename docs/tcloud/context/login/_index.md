---
linkTitle: "tcloud context login"
title: "context login"
slug: tcloud_context_login
url: /docs/tcloud/context/login/
weight: 9980
cascade:
  type: docs
---
## tcloud context login

Login to Thalassa Cloud

### Synopsis

Login to Thalassa Cloud using browser OIDC (default when no credentials are provided), a personal access token, access token, or OIDC client credentials for the current context.

```
tcloud context login [flags]
```

### Options

```
      --browser       log in through the browser using OIDC (default when no other credentials are provided)
  -h, --help          help for login
      --name string   name of the context (default "default")
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

