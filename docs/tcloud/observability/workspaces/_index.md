---
linkTitle: "tcloud observability workspaces"
title: "observability workspaces"
slug: tcloud_observability_workspaces
url: /docs/tcloud/observability/workspaces/
weight: 9714
cascade:
  type: docs
---
## tcloud observability workspaces

Manage observability workspaces

### Synopsis

Create and manage unified Prometheus + Loki workspaces.

Each workspace exposes remote-write / push URLs for ingest and query URLs for
Prometheus and Loki in the current organisation/project context.

### Options

```
  -h, --help   help for workspaces
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

* [tcloud observability](/docs/tcloud/tcloud_observability/)	 - Manage observability workspaces (metrics and logs)
* [tcloud observability workspaces create](/docs/tcloud/observability/workspaces_create/)	 - Create an observability workspace
* [tcloud observability workspaces delete](/docs/tcloud/observability/workspaces_delete/)	 - Delete an observability workspace
* [tcloud observability workspaces get](/docs/tcloud/observability/workspaces_get/)	 - Show an observability workspace
* [tcloud observability workspaces list](/docs/tcloud/observability/workspaces_list/)	 - List observability workspaces
* [tcloud observability workspaces update](/docs/tcloud/observability/workspaces_update/)	 - Update an observability workspace (only set flags are changed)

