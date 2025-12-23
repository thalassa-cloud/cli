---
linkTitle: "tcloud compute machines stop"
title: "compute machines stop"
slug: tcloud_compute_machines_stop
url: /docs/tcloud/compute/machines_stop/
weight: 9992
cascade:
  type: docs
---
## tcloud compute machines stop

Stop a machine

### Synopsis

Stop a machine to stop it from running. This command will stop the machine and all the services associated with it.

```
tcloud compute machines stop [flags]
```

### Options

```
  -h, --help   help for stop
  -w, --wait   Wait for the machine to be stopped
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

