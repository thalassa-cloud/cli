---
linkTitle: "tcloud projects view"
title: "projects view"
slug: tcloud_projects_view
url: /docs/tcloud/projects/view/
weight: 9730
cascade:
  type: docs
---
## tcloud projects view

View a project

```
tcloud projects view <project> [flags]
```

### Options

```
      --exact-time   Show full timestamps instead of relative time
  -h, --help         help for view
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

* [tcloud projects](/docs/tcloud/tcloud_projects/)	 - Manage projects (beta)

