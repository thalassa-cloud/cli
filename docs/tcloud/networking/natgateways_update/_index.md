---
linkTitle: "tcloud networking natgateways update"
title: "networking natgateways update"
slug: tcloud_networking_natgateways_update
url: /docs/tcloud/networking/natgateways_update/
weight: 9799
cascade:
  type: docs
---
## tcloud networking natgateways update

Update a NAT gateway

### Synopsis

Update properties of an existing NAT gateway. Unspecified fields are preserved.

```
tcloud networking natgateways update [flags]
```

### Examples

```
tcloud networking natgateways update ngw-123 --name egress-prod
tcloud networking natgateways update ngw-123 --reserved-ip rip-456
tcloud networking natgateways update ngw-123 --detach-reserved-ip
```

### Options

```
      --annotations strings       Annotations in key=value format
      --description string        Description of the NAT gateway
      --detach-reserved-ip        Detach the currently associated reserved IP
  -h, --help                      help for update
      --labels strings            Labels in key=value format
      --name string               Name of the NAT gateway
      --reserved-ip string        Reserved IP identity to attach or replace
      --security-groups strings   Security group identities to attach
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

