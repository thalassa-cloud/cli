---
linkTitle: "tcloud networking natgateways create"
title: "networking natgateways create"
slug: tcloud_networking_natgateways_create
url: /docs/tcloud/networking/natgateways_create/
weight: 9802
cascade:
  type: docs
---
## tcloud networking natgateways create

Create a NAT gateway

### Synopsis

Create a new NAT gateway in the specified subnet.

```
tcloud networking natgateways create [flags]
```

### Examples

```
tcloud networking natgateways create --name egress --subnet subnet-123
tcloud networking natgateways create --name egress --subnet subnet-123 --configure-default-route --wait
```

### Options

```
      --annotations strings       Annotations in key=value format
      --configure-default-route   Configure the default route for the subnet route table
      --description string        Description of the NAT gateway
  -h, --help                      help for create
      --labels strings            Labels in key=value format
      --name string               Name of the NAT gateway
      --reserved-ip string        Reserved IP identity to attach
      --security-groups strings   Security group identities to attach
      --subnet string             Subnet identity, slug, or name
      --wait                      Wait for the NAT gateway to have an endpoint
      --wait-timeout duration     Maximum time to wait for the NAT gateway endpoint (default 20m0s)
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

* [tcloud networking natgateways](/docs/tcloud/networking/natgateways/)	 - Manage NAT gateways

