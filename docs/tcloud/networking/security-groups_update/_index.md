---
linkTitle: "tcloud networking security-groups update"
title: "networking security-groups update"
slug: tcloud_networking_security-groups_update
url: /docs/tcloud/networking/security-groups_update/
weight: 9769
cascade:
  type: docs
---
## tcloud networking security-groups update

Update security group metadata

### Synopsis

Update name, description, labels, annotations, or allow-same-group. Rules are left unchanged.

```
tcloud networking security-groups update SECURITY_GROUP [flags]
```

### Options

```
      --allow-same-group      Allow traffic between instances in the same security group
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Description
  -h, --help                  help for update
      --labels strings        Labels as key=value (repeatable)
      --name string           Name
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

* [tcloud networking security-groups](/docs/tcloud/networking/security-groups/)	 - Manage security groups

