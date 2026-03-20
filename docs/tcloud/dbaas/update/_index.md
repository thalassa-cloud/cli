---
linkTitle: "tcloud dbaas update"
title: "dbaas update"
slug: tcloud_dbaas_update
url: /docs/tcloud/dbaas/update/
weight: 9960
cascade:
  type: docs
---
## tcloud dbaas update

Update a database cluster

### Synopsis

Update properties of an existing database cluster.

```
tcloud dbaas update [flags]
```

### Options

```
      --annotations strings    Annotations in key=value format (can be specified multiple times)
      --delete-protection      Enable or disable delete protection
      --description string     Description of the database cluster
  -h, --help                   help for update
      --instance-type string   Instance type
      --labels strings         Labels in key=value format (can be specified multiple times)
      --name string            Name of the database cluster
      --no-header              Do not print the header
      --replicas int           Number of replicas (default -1)
      --storage int            Storage size in GB
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

