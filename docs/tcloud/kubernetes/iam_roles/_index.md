---
linkTitle: "tcloud kubernetes iam roles"
title: "kubernetes iam roles"
slug: tcloud_kubernetes_iam_roles
url: /docs/tcloud/kubernetes/iam_roles/
weight: 9834
cascade:
  type: docs
---
## tcloud kubernetes iam roles

Kubernetes cluster IAM roles, permission rules, and bindings

### Synopsis

Custom Kubernetes cluster roles define RBAC-style permission rules and can be bound to users,
teams, or service accounts. System roles may be read-only; the API enforces what you can change.

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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud kubernetes iam](/docs/tcloud/kubernetes/iam/)	 - Kubernetes cluster IAM roles and bindings
* [tcloud kubernetes iam roles bindings](/docs/tcloud/kubernetes/iam_roles_bindings/)	 - Kubernetes cluster role bindings (who receives the role)
* [tcloud kubernetes iam roles create](/docs/tcloud/kubernetes/iam_roles_create/)	 - Create a custom Kubernetes cluster IAM role
* [tcloud kubernetes iam roles delete](/docs/tcloud/kubernetes/iam_roles_delete/)	 - Delete a custom Kubernetes cluster IAM role
* [tcloud kubernetes iam roles get](/docs/tcloud/kubernetes/iam_roles_get/)	 - Show a Kubernetes cluster role including rules and bindings
* [tcloud kubernetes iam roles list](/docs/tcloud/kubernetes/iam_roles_list/)	 - List Kubernetes cluster IAM roles
* [tcloud kubernetes iam roles rules](/docs/tcloud/kubernetes/iam_roles_rules/)	 - Permission rules on a Kubernetes cluster role

