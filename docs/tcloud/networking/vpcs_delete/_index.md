---
linkTitle: "tcloud networking vpcs delete"
title: "networking vpcs delete"
slug: tcloud_networking_vpcs_delete
url: /docs/tcloud/networking/vpcs_delete/
weight: 9938
cascade:
  type: docs
---
## tcloud networking vpcs delete

Delete VPC(s)

### Synopsis

Delete VPC(s) by identity or label selector. This command will delete the VPC(s) and all associated resources.

```
tcloud networking vpcs delete [flags]
```

### Examples

```
tcloud networking vpcs delete vpc-123
tcloud networking vpcs delete vpc-123 vpc-456 --wait
tcloud networking vpcs delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter VPCs (format: key1=value1,key2=value2)
      --wait              Wait for the VPC(s) to be deleted
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

* [tcloud networking vpcs](/docs/tcloud/networking/vpcs/)	 - Manage VPCs

