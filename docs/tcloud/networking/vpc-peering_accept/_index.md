---
linkTitle: "tcloud networking vpc-peering accept"
title: "networking vpc-peering accept"
slug: tcloud_networking_vpc-peering_accept
url: /docs/tcloud/networking/vpc-peering_accept/
weight: 9939
cascade:
  type: docs
---
## tcloud networking vpc-peering accept

Accept a VPC peering connection

### Synopsis

Accept a pending VPC peering connection request

```
tcloud networking vpc-peering accept [flags]
```

### Examples

```
tcloud networking vpc-peering accept vpcpc-123
tcloud networking vpc-peering accept vpcpc-123 --force
```

### Options

```
      --force       Force the acceptance and skip the confirmation
  -h, --help        help for accept
      --no-header   Do not print the header
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

