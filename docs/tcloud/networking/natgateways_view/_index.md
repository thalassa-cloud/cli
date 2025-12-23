---
linkTitle: "tcloud networking natgateways view"
title: "networking natgateways view"
slug: tcloud_networking_natgateways_view
url: /docs/tcloud/networking/natgateways_view/
weight: 9962
cascade:
  type: docs
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

