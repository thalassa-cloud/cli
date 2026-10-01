---
linkTitle: "tcloud dir show"
title: "dir show"
slug: tcloud_dir_show
url: /docs/tcloud/dir/show/
weight: 9939
cascade:
  type: docs
---
## tcloud dir show

Show the directory-local config in effect

### Synopsis

Show the nearest .thalassa (or .thalassa/config.yaml) found by walking up from the current directory.

```
tcloud dir show [flags]
```

### Options

```
  -h, --help   help for show
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

* [tcloud dir](/docs/tcloud/tcloud_dir/)	 - Manage directory-local Thalassa defaults

