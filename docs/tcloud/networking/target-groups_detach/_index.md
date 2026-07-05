---
linkTitle: "tcloud networking target-groups detach"
title: "networking target-groups detach"
slug: tcloud_networking_target-groups_detach
url: /docs/tcloud/networking/target-groups_detach/
weight: 9841
cascade:
  type: docs
---
## tcloud networking target-groups detach

Detach a target from a target group

### Synopsis

Detach a target attachment from a target group.

```
tcloud networking target-groups detach TARGET_GROUP ATTACHMENT [flags]
```

### Examples

```
tcloud networking target-groups detach tg-123 attachment-456
```

### Options

```
  -h, --help   help for detach
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

