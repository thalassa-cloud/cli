---
linkTitle: "tcloud networking loadbalancers list"
title: "networking loadbalancers list"
slug: tcloud_networking_loadbalancers_list
url: /docs/tcloud/networking/loadbalancers_list/
weight: 9869
cascade:
  type: docs
---
## tcloud networking loadbalancers list

List load balancers

### Synopsis

List load balancers within your organisation.

```
tcloud networking loadbalancers list [flags]
```

### Examples

```
tcloud networking loadbalancers list
tcloud networking loadbalancers list --region us-west-1
tcloud networking loadbalancers list --vpc vpc-123 --no-header
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Filter by region
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

* [tcloud networking loadbalancers](/docs/tcloud/networking/loadbalancers/)	 - Manage load balancers

