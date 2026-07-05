---
linkTitle: "tcloud registry namespaces retention run"
title: "registry namespaces retention run"
slug: tcloud_registry_namespaces_retention_run
url: /docs/tcloud/registry/namespaces_retention_run/
weight: 9800
cascade:
  type: docs
---
## tcloud registry namespaces retention run

Run the retention policy for a namespace

```
tcloud registry namespaces retention run [flags]
```

### Examples

```
tcloud registry namespaces retention run --namespace crns-123
```

### Options

```
  -h, --help               help for run
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

* [tcloud registry namespaces retention](/docs/tcloud/registry/namespaces_retention/)	 - Run retention policies

