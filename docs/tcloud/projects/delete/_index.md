---
linkTitle: "tcloud projects delete"
title: "projects delete"
slug: tcloud_projects_delete
url: /docs/tcloud/projects/delete/
weight: 9733
cascade:
  type: docs
---
## tcloud projects delete

Delete a project

```
tcloud projects delete <project> [flags]
```

### Options

```
      --force   Skip the confirmation prompt and delete
  -h, --help    help for delete
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

