---
linkTitle: "tcloud networking vpc-peering delete"
title: "networking vpc-peering delete"
slug: tcloud_networking_vpc-peering_delete
url: /docs/tcloud/networking/vpc-peering_delete/
weight: 9937
cascade:
  type: docs
---
## tcloud networking vpc-peering delete

Delete VPC peering connection(s)

### Synopsis

Delete VPC peering connection(s) by identity or label selector

```
tcloud networking vpc-peering delete [flags]
```

### Examples

```
tcloud networking vpc-peering delete vpcpc-123
tcloud networking vpc-peering delete vpcpc-123 vpcpc-456 --force
tcloud networking vpc-peering delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter connections (format: key1=value1,key2=value2)
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

