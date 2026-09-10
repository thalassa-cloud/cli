---
linkTitle: "tcloud projects list"
title: "projects list"
slug: tcloud_projects_list
url: /docs/tcloud/projects/list/
weight: 9732
cascade:
  type: docs
---
## tcloud projects list

List projects in the current organisation

```
tcloud projects list [flags]
```

### Options

```
  -h, --help           help for list
      --include-root   include organisation root (no project) as the first entry
      --no-header      do not print headers
      --slug-only      only print the slug
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

