---
linkTitle: "tcloud observability"
title: "observability"
slug: tcloud_observability
url: /docs/tcloud/tcloud_observability/
weight: 9713
cascade:
  type: docs
---
## tcloud observability

Manage observability workspaces (metrics and logs)

### Synopsis

Manage Thalassa observability workspaces that provide Prometheus (Cortex) metrics
and Loki log ingestion/query endpoints for the current organisation/project scope.

### Options

```
  -h, --help   help for observability
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud observability workspaces](/docs/tcloud/observability/workspaces/)	 - Manage observability workspaces

