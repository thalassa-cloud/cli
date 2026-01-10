---
linkTitle: "tcloud networking vpc-peering update"
title: "networking vpc-peering update"
slug: tcloud_networking_vpc-peering_update
url: /docs/tcloud/networking/vpc-peering_update/
weight: 9934
cascade:
  type: docs
---
## tcloud networking vpc-peering update

Update a VPC peering connection

```
tcloud networking vpc-peering update [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description of the VPC peering connection
  -h, --help                  help for update
      --labels strings        Labels in key=value format
      --name string           Name of the VPC peering connection
      --no-header             Do not print the header
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

