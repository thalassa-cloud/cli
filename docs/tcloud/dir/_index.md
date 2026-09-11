---
linkTitle: "tcloud dir"
title: "dir"
slug: tcloud_dir
url: /docs/tcloud/tcloud_dir/
weight: 9974
cascade:
  type: docs
---
## tcloud dir

Manage directory-local Thalassa defaults

### Synopsis

Manage optional project-local defaults discovered from a `.thalassa` file (or `.thalassa/config.yaml`) by walking up from the current directory.

Directory config can set a preferred CLI context, organisation, project, and Kubernetes cluster.

```yaml
context: prod
organisation: acme
project: platform
kubernetes:
  cluster: prod-cluster
```

Disable with `--ignore-dir-config` or `THALASSA_DIR_CONFIG=0`.

### Options

```
  -h, --help   help for dir
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud dir init](/docs/tcloud/dir/init/)	 - Write a .thalassa file in the current directory
* [tcloud dir show](/docs/tcloud/dir/show/)	 - Show the directory-local config in effect
