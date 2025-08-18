---
date: 2025-08-14T00:38:24+02:00
linkTitle: "tcloud dbaas versions"
title: "dbaas versions"
slug: tcloud_dbaas_versions
url: /docs/tcloud/tcloud_dbaas_versions/
weight: 9978
---
## tcloud dbaas versions

Get a list of database engine versions

### Synopsis

Get a list of available database engine versions for a specific engine

```
tcloud dbaas versions [flags]
```

### Examples

```
tcloud dbaas versions --engine postgres
tcloud dbaas versions --engine postgres --no-header
```

### Options

```
      --engine string   Database engine type (e.g., postgres)
  -h, --help            help for versions
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

