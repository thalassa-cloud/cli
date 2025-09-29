---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud networking security-groups view"
title: "networking security-groups view"
slug: tcloud_networking_security-groups_view
url: /docs/tcloud/tcloud_networking_security-groups_view/
weight: 9958
---
## tcloud networking security-groups view

View security group details

### Synopsis

View detailed information about a specific security group

```
tcloud networking security-groups view [flags]
```

### Examples

```
tcloud networking security-groups view sg-123
tcloud networking security-groups view sg-123 --output yaml
```

### Options

```
  -h, --help            help for view
  -o, --output string   Output format (yaml)
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

