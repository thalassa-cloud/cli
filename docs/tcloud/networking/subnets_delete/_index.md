---
linkTitle: "tcloud networking subnets delete"
title: "networking subnets delete"
slug: tcloud_networking_subnets_delete
url: /docs/tcloud/networking/subnets_delete/
weight: 9942
cascade:
  type: docs
---
## tcloud networking subnets delete

Delete subnet(s)

### Synopsis

Delete subnet(s) by identity or label selector.

```
tcloud networking subnets delete [flags]
```

### Examples

```
tcloud networking subnets delete subnet-123
tcloud networking subnets delete subnet-123 subnet-456 --wait
tcloud networking subnets delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter subnets (format: key1=value1,key2=value2)
      --wait              Wait for the subnet(s) to be deleted
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

* [tcloud networking subnets](/docs/tcloud/networking/subnets/)	 - Manage subnets

