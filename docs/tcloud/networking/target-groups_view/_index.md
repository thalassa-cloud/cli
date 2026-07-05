---
linkTitle: "tcloud networking target-groups view"
title: "networking target-groups view"
slug: tcloud_networking_target-groups_view
url: /docs/tcloud/networking/target-groups_view/
weight: 9837
cascade:
  type: docs
---
## tcloud networking target-groups view

View target group details

### Synopsis

View detailed information about a target group.

```
tcloud networking target-groups view TARGET_GROUP [flags]
```

### Examples

```
tcloud networking target-groups view tg-123
tcloud networking target-groups view tg-123 --output yaml
```

### Options

```
  -h, --help            help for view
  -o, --output string   Output format (yaml)
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

