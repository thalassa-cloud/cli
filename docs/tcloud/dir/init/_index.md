---
linkTitle: "tcloud dir init"
title: "dir init"
slug: tcloud_dir_init
url: /docs/tcloud/dir/init/
weight: 9972
cascade:
  type: docs
---
## tcloud dir init

Write a .thalassa file in the current directory

### Synopsis

Write a `.thalassa` file in the current directory from the effective CLI context and an optional Kubernetes cluster.

The file contains names and references only. Credentials stay in `~/.tcloud` or the system credential store. Existing files are not overwritten unless `--force` is set.

```
tcloud dir init [flags]
```

### Examples

```
tcloud dir init --cluster prod-cluster
```

### Options

```
      --cluster string   Kubernetes cluster identity, name, or slug to pin in the file
      --force            Overwrite an existing .thalassa file
  -h, --help             help for init
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
