---
linkTitle: "tcloud projects update"
title: "projects update"
slug: tcloud_projects_update
url: /docs/tcloud/projects/update/
weight: 9731
cascade:
  type: docs
---
## tcloud projects update

Update a project (only flags you set are changed)

```
tcloud projects update <project> [flags]
```

### Options

```
      --annotations strings   Replace annotations (key=value, repeatable)
      --description string    Project description
  -h, --help                  help for update
      --labels strings        Replace labels (key=value, repeatable)
      --name string           Project display name
      --no-header             do not print headers
      --parent string         Parent project identity or slug (empty clears parent)
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

