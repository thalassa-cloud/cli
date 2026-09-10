---
linkTitle: "tcloud registry namespaces list"
title: "registry namespaces list"
slug: tcloud_registry_namespaces_list
url: /docs/tcloud/registry/namespaces_list/
weight: 9715
cascade:
  type: docs
---
## tcloud registry namespaces list

List container registry namespaces

```
tcloud registry namespaces list [flags]
```

### Options

```
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Filter by region
  -l, --selector string   Label selector (format: key1=value1,key2=value2)
      --show-labels       Show labels
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

* [tcloud registry namespaces](/docs/tcloud/registry/namespaces/)	 - Manage container registry namespaces

