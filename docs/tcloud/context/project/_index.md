---
linkTitle: "tcloud context project"
title: "context project"
slug: tcloud_context_project
url: /docs/tcloud/context/project/
weight: 9978
cascade:
  type: docs
---
## tcloud context project

Set the project in the current context

### Synopsis

Set the project in the current context. Accepts a project identity or slug; the identity is stored in the context. Use "root" to clear the project and scope commands to the organisation root.

```
tcloud context project <project> [flags]
```

### Examples

```
tcloud context project <project>
```

### Options

```
  -h, --help   help for project
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

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

