---
linkTitle: "tcloud projects"
title: "projects"
slug: tcloud_projects
url: /docs/tcloud/tcloud_projects/
weight: 9815
cascade:
  type: docs
---
## tcloud projects

Manage projects (private beta)

### Synopsis

Manage projects within the organisation selected in your context.

Note: This command is in private beta and requires the project feature gate to be enabled on your organisation.

### Options

```
  -h, --help   help for projects
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud projects list](/docs/tcloud/projects/list/)	 - List projects in the current organisation

