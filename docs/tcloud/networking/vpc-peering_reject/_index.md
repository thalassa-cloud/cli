---
linkTitle: "tcloud networking vpc-peering reject"
title: "networking vpc-peering reject"
slug: tcloud_networking_vpc-peering_reject
url: /docs/tcloud/networking/vpc-peering_reject/
weight: 9935
cascade:
  type: docs
---
## tcloud networking vpc-peering reject

Reject a VPC peering connection

### Synopsis

Reject a pending VPC peering connection request

```
tcloud networking vpc-peering reject [flags]
```

### Examples

```
tcloud networking vpc-peering reject vpcpc-123
tcloud networking vpc-peering reject vpcpc-123 --reason 'Not needed'
tcloud networking vpc-peering reject vpcpc-123 --force
```

### Options

```
      --force           Force the rejection and skip the confirmation
  -h, --help            help for reject
      --no-header       Do not print the header
      --reason string   Reason for rejecting the peering connection
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

* [tcloud networking vpc-peering](/docs/tcloud/networking/vpc-peering/)	 - Manage VPC peering connections

