---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud networking security-groups list"
title: "networking security-groups list"
slug: tcloud_networking_security-groups_list
url: /docs/tcloud/tcloud_networking_security-groups_list/
weight: 9959
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
      --exact-time   Show exact time instead of relative time
  -h, --help         help for list
      --no-header    Do not print the header
      --vpc string   Filter by VPC
```

### Options inherited from parent commands

```
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud networking security-groups](/docs/tcloud/tcloud_networking_security-groups/)	 - Manage security groups

