---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud context use"
title: "context use"
slug: tcloud_context_use
url: /docs/tcloud/tcloud_context_use/
weight: 9983
---
## tcloud context use

Set the current context

### Synopsis

Set the current context (or the context set with the --context flag)

```
tcloud context use <context> [flags]
```

### Examples

```
tcloud context use <context>
```

### Options

```
  -h, --help   help for use
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

