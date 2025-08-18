---
date: 2025-08-14T00:38:24+02:00
linkTitle: "tcloud context delete"
title: "context delete"
slug: tcloud_context_delete
url: /docs/tcloud/tcloud_context_delete/
weight: 9989
---
## tcloud context delete

Delete a context

### Synopsis

Delete a context from the config

```
tcloud context delete [flags]
```

### Examples

```
tcloud context delete <context>
```

### Options

```
  -h, --help   help for delete
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

