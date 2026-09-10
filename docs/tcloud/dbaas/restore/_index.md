---
linkTitle: "tcloud dbaas restore"
title: "dbaas restore"
slug: tcloud_dbaas_restore
url: /docs/tcloud/dbaas/restore/
weight: 9945
cascade:
  type: docs
---
## tcloud dbaas restore

Restore a database cluster from a backup

### Synopsis

Create a new database cluster from an existing backup.

Restore always provisions a new cluster. Optionally specify a point-in-time
recovery target (timestamp or LSN). In a terminal with fzf installed, missing
cluster, backup, and subnet values can be selected interactively.

```
tcloud dbaas restore [flags]
```

### Examples

```
tcloud dbaas restore --from-backup bak-123 --name restored-db --subnet subnet-123
tcloud dbaas restore --cluster dbc-123 --target-time 2023-12-25T10:00:00Z --name restored-db
```

### Options

```
      --annotations strings     Annotations in key=value format (can be specified multiple times)
      --backup-store string     Backup store identity to attach to the restored cluster
      --cluster string          Source database cluster identity
      --delete-protection       Enable delete protection
      --description string      Description of the restored database cluster
      --force                   Skip confirmation and allow restore from backups with incomplete WAL coverage
      --from-backup string      Backup identity to restore from
  -h, --help                    help for restore
      --instance-type string    Instance type (defaults to the source cluster)
      --labels strings          Labels in key=value format (can be specified multiple times)
      --name string             Name of the restored database cluster (required)
      --no-header               Do not print the header
      --replicas int            Number of replicas (defaults to the source cluster) (default -1)
      --storage int             Storage size in GB (defaults to the source cluster)
      --subnet string           Subnet identity, slug, or name (defaults to the source cluster)
      --target-lsn string       Point-in-time recovery target LSN
      --target-time string      Timestamp to restore to. Accepts RFC3339 (e.g. 2023-12-25T10:00:00Z) or barman format (YYYY-MM-DD HH:MM:SS.00000±TZ, e.g. '2023-08-11 11:14:21.00000+02')
      --volume-type string      Volume type (default "block")
      --vpc string              VPC identity, slug, or name
      --wait                    Wait for the restored cluster to be available before returning
      --wait-timeout duration   Maximum time to wait for the restored cluster to be available (default 20m0s)
      --with-backup-bucket      Provision a backup store for the restored cluster
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

* [tcloud dbaas](/docs/tcloud/tcloud_dbaas/)	 - Manage database clusters and related services

