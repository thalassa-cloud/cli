---
linkTitle: "tcloud iam roles create"
title: "iam roles create"
slug: tcloud_iam_roles_create
url: /docs/tcloud/iam/roles_create/
weight: 9934
cascade:
  type: docs
---
## tcloud iam roles create

Create a custom organisation role

```
tcloud iam roles create [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Role description
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --name string           Role name
      --no-header             Do not print table headers
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

* [tcloud iam roles](/docs/tcloud/iam/roles/)	 - Organisation IAM roles, permission rules, and bindings

