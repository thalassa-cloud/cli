---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud context current"
title: "context current"
slug: tcloud_context_current
url: /docs/tcloud/tcloud_context_current/
weight: 9990
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

