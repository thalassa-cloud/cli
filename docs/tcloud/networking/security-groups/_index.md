---
linkTitle: "tcloud networking security-groups"
title: "networking security-groups"
slug: tcloud_networking_security-groups
url: /docs/tcloud/networking/security-groups/
weight: 9954
cascade:
  type: docs
---
## tcloud networking security-groups

Manage security groups

### Synopsis

Manage security groups and their rules within the Thalassa Cloud Platform

### Examples

```
tcloud networking security-groups list
tcloud networking security-groups create --name my-sg --vpc vpc-123
tcloud networking security-groups delete sg-456
```

### Options

```
  -h, --help   help for security-groups
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

* [tcloud networking](/docs/tcloud/tcloud_networking/)	 - Manage networking resources
* [tcloud networking security-groups create](/docs/tcloud/networking/security-groups_create/)	 - Create a security group
* [tcloud networking security-groups delete](/docs/tcloud/networking/security-groups_delete/)	 - Delete security group(s)
* [tcloud networking security-groups list](/docs/tcloud/networking/security-groups_list/)	 - Get a list of security groups
* [tcloud networking security-groups view](/docs/tcloud/networking/security-groups_view/)	 - View security group details

