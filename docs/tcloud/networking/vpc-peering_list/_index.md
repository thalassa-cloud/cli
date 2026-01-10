---
linkTitle: "tcloud networking vpc-peering list"
title: "networking vpc-peering list"
slug: tcloud_networking_vpc-peering_list
url: /docs/tcloud/networking/vpc-peering_list/
weight: 9936
cascade:
  type: docs
---
## tcloud networking vpc-peering list

List VPC peering connections

```
tcloud networking vpc-peering list [flags]
```

### Options

```
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter connections (format: key1=value1,key2=value2)
      --show-labels       Show labels
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

