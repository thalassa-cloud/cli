---
linkTitle: "tcloud networking security-groups create"
title: "networking security-groups create"
slug: tcloud_networking_security-groups_create
url: /docs/tcloud/networking/security-groups_create/
weight: 9958
cascade:
  type: docs
---
## tcloud networking security-groups create

Create a security group

### Synopsis

Create a new security group within your organisation

```
tcloud networking security-groups create [flags]
```

### Examples

```
tcloud networking security-groups create --name my-sg --vpc vpc-123
tcloud networking security-groups create --name my-sg --vpc vpc-123 --description 'My security group' --allow-same-group
```

### Options

```
      --allow-same-group     Allow traffic between instances in the same security group
      --description string   Description of the security group
  -h, --help                 help for create
      --name string          Name of the security group
      --vpc string           VPC identity where the security group will be created
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

