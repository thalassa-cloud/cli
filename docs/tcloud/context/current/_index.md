---
linkTitle: "tcloud context current"
title: "context current"
slug: tcloud_context_current
url: /docs/tcloud/context/current/
weight: 9988
cascade:
  type: docs
---
## tcloud context current

Shows the current context

### Synopsis

Shows the current context (or the context set with the --context flag)

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
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

