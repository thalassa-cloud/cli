---
linkTitle: "tcloud networking reserved-ips update"
title: "networking reserved-ips update"
slug: tcloud_networking_reserved-ips_update
url: /docs/tcloud/networking/reserved-ips_update/
weight: 9791
cascade:
  type: docs
---
## tcloud networking reserved-ips update

Update a reserved IP address

### Synopsis

Update metadata of an existing reserved IP address. Unspecified fields are preserved from the current resource.

```
tcloud networking reserved-ips update [flags]
```

### Examples

```
tcloud networking reserved-ips update rip-123 --name new-name
tcloud networking reserved-ips update rip-123 --description 'updated'
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description of the reserved IP
  -h, --help                  help for update
      --labels strings        Labels in key=value format
      --name string           Name of the reserved IP
      --no-header             Do not print the header
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

* [tcloud networking reserved-ips](/docs/tcloud/networking/reserved-ips/)	 - Manage reserved IP addresses

