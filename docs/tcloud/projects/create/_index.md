---
linkTitle: "tcloud projects create"
title: "projects create"
slug: tcloud_projects_create
url: /docs/tcloud/projects/create/
weight: 9734
cascade:
  type: docs
---
## tcloud projects create

Create a project in the current organisation

```
tcloud projects create [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Project description
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --name string           Project display name
      --no-header             do not print headers
      --parent string         Parent project identity or slug
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

