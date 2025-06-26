---
date: 2025-06-26T11:24:44+02:00
linkTitle: "tcloud networking vpcs create"
title: "networking vpcs create"
slug: tcloud_networking_vpcs_create
url: /docs/tcloud/tcloud_networking_vpcs_create/
weight: 9964
---
## tcloud networking vpcs create

Create a vpc

```
tcloud networking vpcs create [flags]
```

### Options

```
      --cidrs strings        CIDRs of the vpc (default [10.0.0.0/16])
      --description string   Description of the vpc
  -h, --help                 help for create
      --name string          Name of the vpc
      --no-header            Do not print the header
      --region string        Region of the vpc
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

* [tcloud networking vpcs](/docs/tcloud/tcloud_networking_vpcs/)	 - Manage VPCs

