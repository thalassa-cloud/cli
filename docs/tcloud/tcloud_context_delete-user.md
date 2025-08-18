---
date: 2025-08-14T00:38:24+02:00
linkTitle: "tcloud context delete-user"
title: "context delete-user"
slug: tcloud_context_delete-user
url: /docs/tcloud/tcloud_context_delete-user/
weight: 9987
---
## tcloud context delete-user

Delete a user

### Synopsis

Delete a user from the config

```
tcloud context delete-user [flags]
```

### Examples

```
tcloud context delete-user <user>
```

### Options

```
  -h, --help   help for delete-user
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

