---
linkTitle: "tcloud networking target-groups list"
title: "networking target-groups list"
slug: tcloud_networking_target-groups_list
url: /docs/tcloud/networking/target-groups_list/
weight: 9758
cascade:
  type: docs
---
## tcloud networking target-groups list

List target groups

### Synopsis

List load balancer target groups within your organisation.

```
tcloud networking target-groups list [flags]
```

### Examples

```
tcloud networking target-groups list
tcloud networking target-groups list --vpc vpc-123
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector (format: key1=value1,key2=value2)
      --show-labels       Show labels
      --vpc string        Filter by VPC
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

