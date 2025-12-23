---
linkTitle: "tcloud compute machines list"
title: "compute machines list"
slug: tcloud_compute_machines_list
url: /docs/tcloud/compute/machines_list/
weight: 9994
cascade:
  type: docs
---
## tcloud compute machines list

Get a list of machines

```
tcloud compute machines list [flags]
```

### Options

```
  -h, --help              help for list
      --no-header         Do not print the header
  -o, --output string     Output format. One of: wide
  -l, --selector string   Label selector to filter machines (format: key1=value1,key2=value2)
      --show-exact-time   Show exact time instead of relative time
      --show-labels       Show labels associated with machines
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

