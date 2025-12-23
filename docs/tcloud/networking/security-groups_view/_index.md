---
linkTitle: "tcloud networking security-groups view"
title: "networking security-groups view"
slug: tcloud_networking_security-groups_view
url: /docs/tcloud/networking/security-groups_view/
weight: 9955
cascade:
  type: docs
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

