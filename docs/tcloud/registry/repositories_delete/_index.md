---
linkTitle: "tcloud registry repositories delete"
title: "registry repositories delete"
slug: tcloud_registry_repositories_delete
url: /docs/tcloud/registry/repositories_delete/
weight: 9709
cascade:
  type: docs
---
## tcloud registry repositories delete

Delete repositories and all artifacts

### Synopsis

Permanently delete repositories and all contained artifacts.

```
tcloud registry repositories delete REPOSITORY [REPOSITORY...] [flags]
```

### Examples

```
tcloud registry repositories delete --namespace crns-123 repo-456 --force
```

### Options

```
      --force              Skip confirmation
  -h, --help               help for delete
      --namespace string   Namespace identity
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

