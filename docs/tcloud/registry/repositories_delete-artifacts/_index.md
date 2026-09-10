---
linkTitle: "tcloud registry repositories delete-artifacts"
title: "registry repositories delete-artifacts"
slug: tcloud_registry_repositories_delete-artifacts
url: /docs/tcloud/registry/repositories_delete-artifacts/
weight: 9708
cascade:
  type: docs
---
## tcloud registry repositories delete-artifacts

Delete artifacts from repositories

### Synopsis

Request deletion of artifacts from repositories without deleting the repository itself.

```
tcloud registry repositories delete-artifacts REPOSITORY [REPOSITORY...] [flags]
```

### Examples

```
tcloud registry repositories delete-artifacts --namespace crns-123 repo-456 --force
```

### Options

```
      --force              Skip confirmation
  -h, --help               help for delete-artifacts
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

