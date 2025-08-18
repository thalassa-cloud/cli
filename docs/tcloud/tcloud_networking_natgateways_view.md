---
date: 2025-08-14T00:38:24+02:00
linkTitle: "tcloud networking natgateways view"
title: "networking natgateways view"
slug: tcloud_networking_natgateways_view
url: /docs/tcloud/tcloud_networking_natgateways_view/
weight: 9965
---
## tcloud networking natgateways view

View NAT gateway details

### Synopsis

View detailed information about a specific NAT gateway

```
tcloud networking natgateways view [flags]
```

### Examples

```
tcloud networking natgateways view ngw-123
tcloud networking natgateways view ngw-123 --output yaml
```

### Options

```
  -h, --help            help for view
  -o, --output string   Output format (yaml)
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

