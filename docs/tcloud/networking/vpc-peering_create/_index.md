---
linkTitle: "tcloud networking vpc-peering create"
title: "networking vpc-peering create"
slug: tcloud_networking_vpc-peering_create
url: /docs/tcloud/networking/vpc-peering_create/
weight: 9938
cascade:
  type: docs
---
## tcloud networking vpc-peering create

Create a VPC peering connection

```
tcloud networking vpc-peering create [flags]
```

### Options

```
      --accepter-organisation string   Identity of the accepter organisation
      --accepter-vpc string            Identity of the accepter VPC
      --annotations strings            Annotations in key=value format
      --auto-accept                    Automatically accept the peering connection (only if requester and accepter are in same region and organisation)
      --description string             Description of the VPC peering connection
  -h, --help                           help for create
      --labels strings                 Labels in key=value format
      --name string                    Name of the VPC peering connection
      --no-header                      Do not print the header
      --requester-vpc string           Identity of the requester VPC
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

