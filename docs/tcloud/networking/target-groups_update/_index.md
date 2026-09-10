---
linkTitle: "tcloud networking target-groups update"
title: "networking target-groups update"
slug: tcloud_networking_target-groups_update
url: /docs/tcloud/networking/target-groups_update/
weight: 9756
cascade:
  type: docs
---
## tcloud networking target-groups update

Update a target group

### Synopsis

Update properties of an existing target group.

```
tcloud networking target-groups update TARGET_GROUP [flags]
```

### Examples

```
tcloud networking target-groups update tg-123 --name web-prod
tcloud networking target-groups update tg-123 --port 8443
```

### Options

```
      --annotations strings           Annotations in key=value format
      --description string            Description of the target group
      --enable-proxy-protocol         Enable proxy protocol
  -h, --help                          help for update
      --labels strings                Labels in key=value format
      --loadbalancing-policy string   Load balancing policy (ROUND_ROBIN, RANDOM, MAGLEV)
      --name string                   Name of the target group
      --port int                      Target port
      --protocol string               Target protocol
      --target-selector strings       Label selector for automatic target assignment (key=value)
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

