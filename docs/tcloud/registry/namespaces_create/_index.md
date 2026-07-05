---
linkTitle: "tcloud registry namespaces create"
title: "registry namespaces create"
slug: tcloud_registry_namespaces_create
url: /docs/tcloud/registry/namespaces_create/
weight: 9803
cascade:
  type: docs
---
## tcloud registry namespaces create

Create a container registry namespace

```
tcloud registry namespaces create [flags]
```

### Examples

```
tcloud registry namespaces create --namespace my-app --region eu-west-1
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description
  -h, --help                  help for create
      --labels strings        Labels in key=value format
      --namespace string      Registry namespace name
      --region string         Region identity, slug, or name
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

* [tcloud registry namespaces](/docs/tcloud/registry/namespaces/)	 - Manage container registry namespaces

