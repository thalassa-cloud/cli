---
linkTitle: "tcloud networking target-groups"
title: "networking target-groups"
slug: tcloud_networking_target-groups
url: /docs/tcloud/networking/target-groups/
weight: 9754
cascade:
  type: docs
---
## tcloud networking target-groups

Manage load balancer target groups

### Synopsis

Manage load balancer target groups, attachments, and backend configuration.

### Examples

```
tcloud networking target-groups list
tcloud networking target-groups create --name web --vpc vpc-123 --port 8080 --protocol http
tcloud networking target-groups attach tg-123 --server machine-123
```

### Options

```
  -h, --help   help for target-groups
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

* [tcloud networking](/docs/tcloud/tcloud_networking/)	 - Manage networking resources
* [tcloud networking target-groups attach](/docs/tcloud/networking/target-groups_attach/)	 - Attach a target to a target group
* [tcloud networking target-groups create](/docs/tcloud/networking/target-groups_create/)	 - Create a target group
* [tcloud networking target-groups delete](/docs/tcloud/networking/target-groups_delete/)	 - Delete target group(s)
* [tcloud networking target-groups detach](/docs/tcloud/networking/target-groups_detach/)	 - Detach a target from a target group
* [tcloud networking target-groups list](/docs/tcloud/networking/target-groups_list/)	 - List target groups
* [tcloud networking target-groups set-attachments](/docs/tcloud/networking/target-groups_set-attachments/)	 - Replace target group attachments
* [tcloud networking target-groups update](/docs/tcloud/networking/target-groups_update/)	 - Update a target group
* [tcloud networking target-groups view](/docs/tcloud/networking/target-groups_view/)	 - View target group details

