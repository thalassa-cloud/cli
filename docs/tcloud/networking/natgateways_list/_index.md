---
linkTitle: "tcloud networking natgateways list"
title: "networking natgateways list"
slug: tcloud_networking_natgateways_list
url: /docs/tcloud/networking/natgateways_list/
weight: 9963
cascade:
  type: docs
---
## tcloud networking natgateways list

Get a list of NAT gateways

### Synopsis

Get a list of NAT gateways within your organisation

```
tcloud networking natgateways list [flags]
```

### Examples

```
tcloud networking natgateways list
tcloud networking natgateways list --region us-west-1
tcloud networking natgateways list --vpc vpc-123 --no-header
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Region of the NAT gateway
  -l, --selector string   Label selector to filter NAT gateways (format: key1=value1,key2=value2)
      --show-labels       Show labels
      --vpc string        VPC of the NAT gateway
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

* [tcloud networking natgateways](/docs/tcloud/networking/natgateways/)	 - Manage NAT gateways

