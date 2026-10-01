---
linkTitle: "tcloud context current"
title: "context current"
slug: tcloud_context_current
url: /docs/tcloud/context/current/
weight: 9986
cascade:
  type: docs
---
## tcloud context current

Shows the current context

### Synopsis

Shows the current context (the --context flag, a directory .thalassa overlay, or current-context in ~/.tcloud)

```
tcloud context current [flags]
```

### Examples

```
tcloud context current
```

### Options

```
  -h, --help   help for current
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

