---
date: 2025-08-14T00:09:06+02:00
linkTitle: "tcloud dbaas"
title: "dbaas"
slug: tcloud_dbaas
url: /docs/tcloud/tcloud_dbaas/
weight: 9976
---
## tcloud dbaas

Manage database clusters and related services

### Synopsis

DBaaS commands to manage your database clusters and related services within the Thalassa Cloud Platform

### Examples

```
tcloud dbaas list
tcloud dbaas instance-types
tcloud dbaas versions --engine postgres
```

### Options

```
  -h, --help   help for dbaas
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud dbaas engines](/docs/tcloud/tcloud_dbaas_engines/)	 - Get a list of database engines
* [tcloud dbaas instance-types](/docs/tcloud/tcloud_dbaas_instance-types/)	 - Get a list of database instance types
* [tcloud dbaas list](/docs/tcloud/tcloud_dbaas_list/)	 - Get a list of database clusters
* [tcloud dbaas versions](/docs/tcloud/tcloud_dbaas_versions/)	 - Get a list of database engine versions

