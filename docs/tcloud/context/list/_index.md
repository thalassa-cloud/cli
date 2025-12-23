---
linkTitle: "tcloud context list"
title: "context list"
slug: tcloud_context_list
url: /docs/tcloud/context/list/
weight: 9984
cascade:
  type: docs
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

