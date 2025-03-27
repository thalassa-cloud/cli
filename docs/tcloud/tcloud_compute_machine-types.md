---
date: 2025-03-27T22:27:48+01:00
linkTitle: "tcloud compute machine-types"
title: "compute machine-types"
slug: tcloud_compute_machine-types
url: /docs/tcloud/tcloud_compute_machine-types/
weight: 9998
---
## tcloud compute machine-types

Get a list of machine types

### Synopsis

Get a list of machine types available in the current organisation

```
tcloud compute machine-types [flags]
```

### Examples

```
thalassa compute machine-types
```

### Options

```
      --category string   Filter by category
  -h, --help              help for machine-types
      --no-header         Do not print the header
  -o, --output string     Output format. One of: wide
      --show-labels       Show labels associated with machines
```

### Options inherited from parent commands

```
      --api string            API endpoint
  -c, --context string        Context name
  -O, --organisation string   Organisation slug or identity
      --token string          Personal access token
```

### SEE ALSO

* [tcloud compute](/docs/tcloud/tcloud_compute/)	 - Manage compute resources

