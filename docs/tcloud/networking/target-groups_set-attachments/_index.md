---
linkTitle: "tcloud networking target-groups set-attachments"
title: "networking target-groups set-attachments"
slug: tcloud_networking_target-groups_set-attachments
url: /docs/tcloud/networking/target-groups_set-attachments/
weight: 9839
cascade:
  type: docs
---
## tcloud networking target-groups set-attachments

Replace target group attachments

### Synopsis

Replace all server and endpoint attachments on a target group. Existing attachments not listed are removed.

```
tcloud networking target-groups set-attachments TARGET_GROUP [flags]
```

### Examples

```
tcloud networking target-groups set-attachments tg-123 --server machine-1 --server machine-2
```

### Options

```
      --endpoint strings   Endpoint identities to attach
  -h, --help               help for set-attachments
      --server strings     Server identities to attach
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

