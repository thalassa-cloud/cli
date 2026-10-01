---
linkTitle: "tcloud observability workspaces create"
title: "observability workspaces create"
slug: tcloud_observability_workspaces_create
url: /docs/tcloud/observability/workspaces_create/
weight: 9719
cascade:
  type: docs
---
## tcloud observability workspaces create

Create an observability workspace

```
tcloud observability workspaces create [flags]
```

### Examples

```
  tcloud observability workspaces create --name prod-metrics --region nl-ams
  tcloud observability workspaces create --name prod-metrics --region nl-ams --retention-days 30 --wait
```

### Options

```
      --annotations strings     Annotations as key=value (repeatable)
      --description string      Workspace description
  -h, --help                    help for create
      --labels strings          Labels as key=value (repeatable)
      --name string             Workspace name (required)
      --region string           Region identity, slug, or name
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

