---
linkTitle: "tcloud networking target-groups delete"
title: "networking target-groups delete"
slug: tcloud_networking_target-groups_delete
url: /docs/tcloud/networking/target-groups_delete/
weight: 9842
cascade:
  type: docs
---
## tcloud networking target-groups delete

Delete target group(s)

### Synopsis

Delete target group(s) by identity or label selector.

```
tcloud networking target-groups delete [TARGET_GROUP...] [flags]
```

### Examples

```
tcloud networking target-groups delete tg-123
tcloud networking target-groups delete --selector env=test --force
```

### Options

```
      --force             Skip confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector (format: key1=value1,key2=value2)
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

