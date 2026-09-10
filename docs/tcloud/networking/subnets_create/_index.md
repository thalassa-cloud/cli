---
linkTitle: "tcloud networking subnets create"
title: "networking subnets create"
slug: tcloud_networking_subnets_create
url: /docs/tcloud/networking/subnets_create/
weight: 9766
cascade:
  type: docs
---
## tcloud networking subnets create

Create a subnet

```
tcloud networking subnets create [flags]
```

### Options

```
      --cidr string             CIDR of the subnet
      --description string      Description of the subnet
  -h, --help                    help for create
      --labels strings          Labels in key=value format (can be specified multiple times)
      --name string             Name of the subnet
      --no-header               Do not print the header
      --vpc string              VPC of the subnet
      --wait                    Wait for the subnet to be ready before returning
      --wait-timeout duration   Maximum time to wait for the subnet to be ready (default 10m0s)
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud networking subnets](/docs/tcloud/networking/subnets/)	 - Manage subnets

