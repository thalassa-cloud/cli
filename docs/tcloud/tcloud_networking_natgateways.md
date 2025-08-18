---
date: 2025-08-14T00:38:24+02:00
linkTitle: "tcloud networking natgateways"
title: "networking natgateways"
slug: tcloud_networking_natgateways
url: /docs/tcloud/tcloud_networking_natgateways/
weight: 9964
---
## tcloud networking natgateways

Manage NAT gateways

### Synopsis

Manage NAT gateways within the Thalassa Cloud Platform. This command will list all the NAT gateways within your organisation.

### Examples

```
tcloud networking natgateways list
tcloud networking natgateways list --region us-west-1
tcloud networking natgateways view ngw-123
```

### Options

```
  -h, --help   help for natgateways
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

* [tcloud networking](/docs/tcloud/tcloud_networking/)	 - Manage networking resources
* [tcloud networking natgateways list](/docs/tcloud/tcloud_networking_natgateways_list/)	 - Get a list of NAT gateways
* [tcloud networking natgateways view](/docs/tcloud/tcloud_networking_natgateways_view/)	 - View NAT gateway details

