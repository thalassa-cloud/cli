---
date: 2025-09-29T22:35:32+02:00
linkTitle: "tcloud dbaas instance-types"
title: "dbaas instance-types"
slug: tcloud_dbaas_instance-types
url: /docs/tcloud/tcloud_dbaas_instance-types/
weight: 9980
---
## tcloud dbaas instance-types

Get a list of database instance types

### Synopsis

Get a list of available database instance types within your organisation

```
tcloud dbaas instance-types [flags]
```

### Examples

```
tcloud dbaas instance-types
tcloud dbaas instance-types --no-header
```

### Options

```
  -h, --help   help for instance-types
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

