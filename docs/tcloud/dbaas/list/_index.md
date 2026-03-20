---
linkTitle: "tcloud dbaas list"
title: "dbaas list"
slug: tcloud_dbaas_list
url: /docs/tcloud/dbaas/list/
weight: 9961
cascade:
  type: docs
---
## tcloud dbaas list

Get a list of database clusters

### Synopsis

Get a list of database clusters within your organisation

```
tcloud dbaas list [flags]
```

### Examples

```
tcloud dbaas list
tcloud dbaas list --no-header
tcloud dbaas list --exact-time
```

### Options

```
      --engine string     Filter by database engine (e.g., postgres)
      --exact-time        Show exact time instead of relative time
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter clusters (format: key1=value1,key2=value2)
      --show-labels       Show labels
      --subnet string     Filter by subnet identity, slug, or name
      --vpc string        Filter by VPC identity, slug, or name
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud dbaas](/docs/tcloud/tcloud_dbaas/)	 - Manage database clusters and related services

