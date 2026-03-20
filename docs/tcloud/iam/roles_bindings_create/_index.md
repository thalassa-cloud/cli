---
linkTitle: "tcloud iam roles bindings create"
title: "iam roles bindings create"
slug: tcloud_iam_roles_bindings_create
url: /docs/tcloud/iam/roles_bindings_create/
weight: 9938
cascade:
  type: docs
---
## tcloud iam roles bindings create

Create a binding to a user, team, or service account

```
tcloud iam roles bindings create <role> [flags]
```

### Options

```
      --annotations strings               Annotations as key=value (repeatable)
      --description string                Binding description
  -h, --help                              help for create
      --labels strings                    Labels as key=value (repeatable)
      --name string                       Binding name
      --no-header                         Do not print table headers
      --scope strings                     Scopes for the binding (repeatable)
      --service-account-identity string   Bind to this service account identity
      --team-identity string              Bind to this team identity
      --user-identity string              Bind to this user identity
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

* [tcloud iam roles bindings](/docs/tcloud/iam/roles_bindings/)	 - Role bindings (who receives the role)

