---
date: 2025-08-14T00:09:06+02:00
linkTitle: "tcloud compute machines delete"
title: "compute machines delete"
slug: tcloud_compute_machines_delete
url: /docs/tcloud/tcloud_compute_machines_delete/
weight: 9997
---
## tcloud compute machines delete

Delete a machine

### Synopsis

Delete a machine. This command will delete the machine and all the services associated with it.

```
tcloud compute machines delete [flags]
```

### Options

```
  -h, --help   help for delete
  -w, --wait   Wait for the machine to be deleted
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

* [tcloud compute machines](/docs/tcloud/tcloud_compute_machines/)	 - Manage virtual machine instances

