---
date: 2025-08-14T00:09:06+02:00
linkTitle: "tcloud compute machines start"
title: "compute machines start"
slug: tcloud_compute_machines_start
url: /docs/tcloud/tcloud_compute_machines_start/
weight: 9995
---
## tcloud compute machines start

Start a machine

### Synopsis

Start a machine to start it from stopped state. This command will start the machine and all the services associated with it.

```
tcloud compute machines start [flags]
```

### Options

```
  -h, --help   help for start
  -w, --wait   Wait for the machine to be started
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

