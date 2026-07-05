---
linkTitle: "tcloud registry repositories"
title: "registry repositories"
slug: tcloud_registry_repositories
url: /docs/tcloud/registry/repositories/
weight: 9791
cascade:
  type: docs
---
## tcloud registry repositories

Manage container registry repositories

### Synopsis

List and manage repositories within a container registry namespace. All commands require --namespace.

### Examples

```
tcloud registry repositories list --namespace crns-123
tcloud registry repositories delete --namespace crns-123 repo-456
```

### Options

```
  -h, --help   help for repositories
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

* [tcloud registry](/docs/tcloud/tcloud_registry/)	 - Manage the Thalassa container registry
* [tcloud registry repositories delete](/docs/tcloud/registry/repositories_delete/)	 - Delete repositories and all artifacts
* [tcloud registry repositories delete-artifacts](/docs/tcloud/registry/repositories_delete-artifacts/)	 - Delete artifacts from repositories
* [tcloud registry repositories list](/docs/tcloud/registry/repositories_list/)	 - List repositories in a namespace
* [tcloud registry repositories view](/docs/tcloud/registry/repositories_view/)	 - View repository details

