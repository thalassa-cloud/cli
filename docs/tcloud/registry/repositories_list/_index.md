---
linkTitle: "tcloud registry repositories list"
title: "registry repositories list"
slug: tcloud_registry_repositories_list
url: /docs/tcloud/registry/repositories_list/
weight: 9793
cascade:
  type: docs
---
## tcloud registry repositories list

List repositories in a namespace

```
tcloud registry repositories list [flags]
```

### Options

```
      --exact-time         Show exact time instead of relative time
  -h, --help               help for list
      --namespace string   Namespace identity
      --no-header          Do not print the header
  -l, --selector string    Label selector (format: key1=value1,key2=value2)
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

* [tcloud registry repositories](/docs/tcloud/registry/repositories/)	 - Manage container registry repositories

