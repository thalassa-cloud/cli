---
linkTitle: "tcloud dbaas create"
title: "dbaas create"
slug: tcloud_dbaas_create
url: /docs/tcloud/dbaas/create/
weight: 9964
cascade:
  type: docs
---
## tcloud dbaas create

Create a database cluster

### Synopsis

Create a new database cluster in the Thalassa Cloud Platform.

```
tcloud dbaas create [flags]
```

### Options

```
      --annotations strings               Annotations in key=value format (can be specified multiple times)
      --backup-object-storage-id string   Backup object storage ID (enables backup storage, requires --with-backup-bucket=false)
      --delete-protection                 Enable delete protection
      --description string                Description of the database cluster
      --engine string                     Database engine (e.g., postgres) (required)
      --engine-version string             Engine version (required)
  -h, --help                              help for create
      --instance-type string              Instance type (required)
      --labels strings                    Labels in key=value format (can be specified multiple times)
      --name string                       Name of the database cluster (required)
      --no-header                         Do not print the header
      --replicas int                      Number of replicas (default: 0)
      --storage int                       Storage size in GB (required)
      --subnet string                     Subnet identity, slug, or name (required)
      --volume-type string                Volume type (default "block")
      --vpc string                        VPC identity, slug, or name
      --wait                              Wait for the database cluster to be available before returning
      --with-backup-bucket                Provision a backup object storage bucket for the database cluster
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

