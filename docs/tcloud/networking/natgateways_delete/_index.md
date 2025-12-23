---
linkTitle: "tcloud networking natgateways delete"
title: "networking natgateways delete"
slug: tcloud_networking_natgateways_delete
url: /docs/tcloud/networking/natgateways_delete/
weight: 9964
cascade:
  type: docs
---
## tcloud networking natgateways delete

Delete NAT gateway(s)

### Synopsis

Delete NAT gateway(s) by identity or label selector. This command will delete the NAT gateway(s) and all associated resources.

```
tcloud networking natgateways delete [flags]
```

### Examples

```
tcloud networking natgateways delete ngw-123
tcloud networking natgateways delete ngw-123 ngw-456 --wait
tcloud networking natgateways delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter NAT gateways (format: key1=value1,key2=value2)
      --wait              Wait for the NAT gateway(s) to be deleted
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

