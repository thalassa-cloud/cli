---
linkTitle: "tcloud networking security-groups delete"
title: "networking security-groups delete"
slug: tcloud_networking_security-groups_delete
url: /docs/tcloud/networking/security-groups_delete/
weight: 9947
cascade:
  type: docs
---
## tcloud networking security-groups delete

Delete security group(s)

### Synopsis

Delete security group(s) by identity or label selector. This command will delete the security group(s) and all their rules.

```
tcloud networking security-groups delete [flags]
```

### Examples

```
tcloud networking security-groups delete sg-123
tcloud networking security-groups delete sg-123 sg-456 --wait
tcloud networking security-groups delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter security groups (format: key1=value1,key2=value2)
      --wait              Wait for the security group(s) to be deleted
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

* [tcloud networking security-groups](/docs/tcloud/networking/security-groups/)	 - Manage security groups

