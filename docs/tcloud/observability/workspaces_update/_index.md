---
linkTitle: "tcloud observability workspaces update"
title: "observability workspaces update"
slug: tcloud_observability_workspaces_update
url: /docs/tcloud/observability/workspaces_update/
weight: 9715
cascade:
  type: docs
---
## tcloud observability workspaces update

Update an observability workspace (only set flags are changed)

```
tcloud observability workspaces update <workspace> [flags]
```

### Options

```
      --annotations strings     Replace annotations (key=value, repeatable)
      --description string      Workspace description
  -h, --help                    help for update
      --labels strings          Replace labels (key=value, repeatable)
      --name string             Workspace name
      --retention-days int      Retention in days for metrics and logs (1-1095)
      --wait                    Wait until the workspace is ready
      --wait-timeout duration   Timeout when waiting with --wait (default 10m0s)
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

* [tcloud observability workspaces](/docs/tcloud/observability/workspaces/)	 - Manage observability workspaces

