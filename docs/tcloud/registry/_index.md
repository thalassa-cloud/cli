---
linkTitle: "tcloud registry"
title: "registry"
slug: tcloud_registry
url: /docs/tcloud/tcloud_registry/
weight: 9790
cascade:
  type: docs
---
## tcloud registry

Manage the Thalassa container registry

### Synopsis

Manage container registry namespaces, repositories, configuration, and retention policies.

### Examples

```
tcloud registry namespaces list
tcloud registry namespaces create --namespace my-app --region eu-west-1
tcloud registry repositories list --namespace crns-123
```

### Options

```
  -h, --help   help for registry
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud registry namespaces](/docs/tcloud/registry/namespaces/)	 - Manage container registry namespaces
* [tcloud registry repositories](/docs/tcloud/registry/repositories/)	 - Manage container registry repositories

