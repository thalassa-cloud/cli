---
linkTitle: "tcloud registry namespaces retention"
title: "registry namespaces retention"
slug: tcloud_registry_namespaces_retention
url: /docs/tcloud/registry/namespaces_retention/
weight: 9799
cascade:
  type: docs
---
## tcloud registry namespaces retention

Run retention policies

### Examples

```
tcloud registry namespaces retention run --namespace crns-123
```

### Options

```
  -h, --help   help for retention
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
* [tcloud registry namespaces retention run](/docs/tcloud/registry/namespaces_retention_run/)	 - Run the retention policy for a namespace

