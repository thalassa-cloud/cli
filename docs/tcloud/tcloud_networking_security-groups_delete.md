---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud networking security-groups delete"
title: "networking security-groups delete"
slug: tcloud_networking_security-groups_delete
url: /docs/tcloud/tcloud_networking_security-groups_delete/
weight: 9960
---
## tcloud networking security-groups delete

Delete a security group

### Synopsis

Delete a security group. This command will delete the security group and all its rules.

```
tcloud networking security-groups delete [flags]
```

### Examples

```
tcloud networking security-groups delete sg-123
tcloud networking security-groups delete sg-123 --wait
```

### Options

```
  -h, --help   help for delete
      --wait   Wait for the security group to be deleted
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

* [tcloud networking security-groups](/docs/tcloud/tcloud_networking_security-groups/)	 - Manage security groups

