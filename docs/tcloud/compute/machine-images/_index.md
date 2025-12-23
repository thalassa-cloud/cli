---
linkTitle: "tcloud compute machine-images"
title: "compute machine-images"
slug: tcloud_compute_machine-images
url: /docs/tcloud/compute/machine-images/
weight: 9997
cascade:
  type: docs
---
## tcloud compute machine-images

Get a list of machine images

### Synopsis

Get a list of machine images available in the current organisation

```
tcloud compute machine-images [flags]
```

### Examples

```
thalassa compute machine-images
```

### Options

```
  -h, --help            help for machine-images
      --no-header       Do not print the header
  -o, --output string   Output format. One of: wide
      --show-labels     Show labels associated with machines
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

* [tcloud compute](/docs/tcloud/tcloud_compute/)	 - Manage compute resources

