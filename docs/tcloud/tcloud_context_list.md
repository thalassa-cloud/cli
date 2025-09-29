---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud context list"
title: "context list"
slug: tcloud_context_list
url: /docs/tcloud/tcloud_context_list/
weight: 9986
---
## tcloud context list

List the contexts

### Synopsis

List the contexts from the config

```
tcloud context list [flags]
```

### Examples

```
tcloud context list
```

### Options

```
  -h, --help           help for list
      --no-header      Do not print the header
      --show-current   Show the current context (default true)
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

