---
linkTitle: "tcloud networking security-groups list"
title: "networking security-groups list"
slug: tcloud_networking_security-groups_list
url: /docs/tcloud/networking/security-groups_list/
weight: 9946
cascade:
  type: docs
---
## tcloud networking security-groups list

Get a list of security groups

### Synopsis

Get a list of security groups within your organisation

```
tcloud networking security-groups list [flags]
```

### Examples

```
tcloud networking security-groups list
tcloud networking security-groups list --no-header
tcloud networking security-groups list --exact-time
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter security groups (format: key1=value1,key2=value2)
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud networking security-groups](/docs/tcloud/networking/security-groups/)	 - Manage security groups

