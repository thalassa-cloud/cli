---
linkTitle: "tcloud compute machines delete"
title: "compute machines delete"
slug: tcloud_compute_machines_delete
url: /docs/tcloud/compute/machines_delete/
weight: 9995
cascade:
  type: docs
---
## tcloud compute machines delete

Delete machine(s)

### Synopsis

Delete machine(s) by identity or label selector. This command will delete the machine(s) and all the services associated with it.

```
tcloud compute machines delete [flags]
```

### Examples

```
tcloud compute machines delete vm-123
tcloud compute machines delete vm-123 vm-456 --wait
tcloud compute machines delete --selector environment=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter machines (format: key1=value1,key2=value2)
  -w, --wait              Wait for the machine(s) to be deleted
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

* [tcloud compute machines](/docs/tcloud/compute/machines/)	 - Manage virtual machine instances

