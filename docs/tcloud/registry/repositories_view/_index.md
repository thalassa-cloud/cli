---
linkTitle: "tcloud registry repositories view"
title: "registry repositories view"
slug: tcloud_registry_repositories_view
url: /docs/tcloud/registry/repositories_view/
weight: 9792
cascade:
  type: docs
---
## tcloud registry repositories view

View repository details

```
tcloud registry repositories view REPOSITORY [flags]
```

### Options

```
  -h, --help               help for view
      --namespace string   Namespace identity
  -o, --output string      Output format (yaml)
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

* [tcloud registry repositories](/docs/tcloud/registry/repositories/)	 - Manage container registry repositories

