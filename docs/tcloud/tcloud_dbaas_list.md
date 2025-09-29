---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud dbaas list"
title: "dbaas list"
slug: tcloud_dbaas_list
url: /docs/tcloud/tcloud_dbaas_list/
weight: 9979
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
      --exact-time   Show exact time instead of relative time
  -h, --help         help for list
      --no-header    Do not print the header
```

### Options inherited from parent commands

```
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

