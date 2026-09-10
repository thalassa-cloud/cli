---
linkTitle: "tcloud registry namespaces delete"
title: "registry namespaces delete"
slug: tcloud_registry_namespaces_delete
url: /docs/tcloud/registry/namespaces_delete/
weight: 9716
cascade:
  type: docs
---
## tcloud registry namespaces delete

Delete container registry namespace(s)

```
tcloud registry namespaces delete [NAMESPACE...] [flags]
```

### Options

```
      --force             Skip confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector (format: key1=value1,key2=value2)
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

