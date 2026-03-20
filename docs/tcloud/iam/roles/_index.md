---
linkTitle: "tcloud iam roles"
title: "iam roles"
slug: tcloud_iam_roles
url: /docs/tcloud/iam/roles/
weight: 9927
cascade:
  type: docs
---
## tcloud iam roles

Organisation IAM roles, permission rules, and bindings

### Synopsis

Custom organisation roles define permission rules and can be bound to users, teams,
or service accounts. System roles may be read-only; the API enforces what you can change.

### Options

```
  -h, --help   help for roles
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

* [tcloud iam](/docs/tcloud/tcloud_iam/)	 - Identity and access management for your organisation
* [tcloud iam roles bindings](/docs/tcloud/iam/roles_bindings/)	 - Role bindings (who receives the role)
* [tcloud iam roles create](/docs/tcloud/iam/roles_create/)	 - Create a custom organisation role
* [tcloud iam roles delete](/docs/tcloud/iam/roles_delete/)	 - Delete a custom organisation role
* [tcloud iam roles get](/docs/tcloud/iam/roles_get/)	 - Show a role including rules and bindings summary
* [tcloud iam roles list](/docs/tcloud/iam/roles_list/)	 - List organisation roles
* [tcloud iam roles rules](/docs/tcloud/iam/roles_rules/)	 - Permission rules on a role

