---
linkTitle: "tcloud networking target-groups attach"
title: "networking target-groups attach"
slug: tcloud_networking_target-groups_attach
url: /docs/tcloud/networking/target-groups_attach/
weight: 9762
cascade:
  type: docs
---
## tcloud networking target-groups attach

Attach a target to a target group

### Synopsis

Attach a server or endpoint to a target group.

```
tcloud networking target-groups attach TARGET_GROUP [flags]
```

### Examples

```
tcloud networking target-groups attach tg-123 --server machine-123
tcloud networking target-groups attach tg-123 --endpoint endpoint-123
```

### Options

```
      --endpoint string   Endpoint identity to attach
  -h, --help              help for attach
      --server string     Server identity to attach
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

* [tcloud networking target-groups](/docs/tcloud/networking/target-groups/)	 - Manage load balancer target groups

