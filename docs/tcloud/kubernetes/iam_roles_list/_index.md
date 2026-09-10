---
linkTitle: "tcloud kubernetes iam roles list"
title: "kubernetes iam roles list"
slug: tcloud_kubernetes_iam_roles_list
url: /docs/tcloud/kubernetes/iam_roles_list/
weight: 9838
cascade:
  type: docs
---
## tcloud kubernetes iam roles list

List Kubernetes cluster IAM roles

```
tcloud kubernetes iam roles list [flags]
```

### Options

```
  -h, --help                    help for list
      --label-selector string   Filter by labels (key=value,key2=value2)
      --no-header               Do not print table headers
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

* [tcloud kubernetes iam roles](/docs/tcloud/kubernetes/iam_roles/)	 - Kubernetes cluster IAM roles, permission rules, and bindings

