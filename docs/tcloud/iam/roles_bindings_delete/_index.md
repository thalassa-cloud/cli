---
linkTitle: "tcloud iam roles bindings delete"
title: "iam roles bindings delete"
slug: tcloud_iam_roles_bindings_delete
url: /docs/tcloud/iam/roles_bindings_delete/
weight: 9935
cascade:
  type: docs
---
## tcloud iam roles bindings delete

Delete a role binding

```
tcloud iam roles bindings delete <role> <binding> [flags]
```

### Options

```
      --force       Skip the confirmation prompt and delete
  -h, --help        help for delete
      --no-header   Do not print table headers
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

* [tcloud iam roles bindings](/docs/tcloud/iam/roles_bindings/)	 - Role bindings (who receives the role)

