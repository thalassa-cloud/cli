---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud networking natgateways list"
title: "networking natgateways list"
slug: tcloud_networking_natgateways_list
url: /docs/tcloud/tcloud_networking_natgateways_list/
weight: 9966
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
      --exact-time      Show exact time instead of relative time
  -h, --help            help for list
      --no-header       Do not print the header
      --region string   Region of the NAT gateway
      --vpc string      VPC of the NAT gateway
```

### Options inherited from parent commands

```
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
  -O, --organisation string    Organisation slug or identity (overrides context)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud networking natgateways](/docs/tcloud/tcloud_networking_natgateways/)	 - Manage NAT gateways

