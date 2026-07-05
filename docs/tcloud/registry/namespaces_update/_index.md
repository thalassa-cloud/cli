---
linkTitle: "tcloud registry namespaces update"
title: "registry namespaces update"
slug: tcloud_registry_namespaces_update
url: /docs/tcloud/registry/namespaces_update/
weight: 9798
cascade:
  type: docs
---
## tcloud registry namespaces update

Update a container registry namespace

```
tcloud registry namespaces update NAMESPACE [flags]
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description
  -h, --help                  help for update
      --labels strings        Labels in key=value format
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

