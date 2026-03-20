---
linkTitle: "tcloud dbaas"
title: "dbaas"
slug: tcloud_dbaas
url: /docs/tcloud/tcloud_dbaas/
weight: 9957
cascade:
  type: docs
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud dbaas backup](/docs/tcloud/dbaas/backup/)	 - Manage database backups
* [tcloud dbaas backup-schedules](/docs/tcloud/dbaas/backup-schedules/)	 - Manage database backup schedules
* [tcloud dbaas create](/docs/tcloud/dbaas/create/)	 - Create a database cluster
* [tcloud dbaas delete](/docs/tcloud/dbaas/delete/)	 - Delete database cluster(s)
* [tcloud dbaas instance-types](/docs/tcloud/dbaas/instance-types/)	 - Get a list of database instance types
* [tcloud dbaas list](/docs/tcloud/dbaas/list/)	 - Get a list of database clusters
* [tcloud dbaas update](/docs/tcloud/dbaas/update/)	 - Update a database cluster
* [tcloud dbaas versions](/docs/tcloud/dbaas/versions/)	 - Get a list of database engine versions
* [tcloud dbaas view](/docs/tcloud/dbaas/view/)	 - View database cluster details

