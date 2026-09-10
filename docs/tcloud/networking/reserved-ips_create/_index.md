---
linkTitle: "tcloud networking reserved-ips create"
title: "networking reserved-ips create"
slug: tcloud_networking_reserved-ips_create
url: /docs/tcloud/networking/reserved-ips_create/
weight: 9795
cascade:
  type: docs
---
## tcloud networking reserved-ips create

Create a reserved IP address

### Synopsis

Create a new reserved public IP address in the specified region.

```
tcloud networking reserved-ips create [flags]
```

### Examples

```
tcloud networking reserved-ips create --name my-ip --region nl-ams
tcloud networking reserved-ips create --name lb-ip --region nl-ams --description 'for production LB'
```

### Options

```
      --annotations strings   Annotations in key=value format
      --description string    Description of the reserved IP
  -h, --help                  help for create
      --labels strings        Labels in key=value format
      --name string           Name of the reserved IP
      --no-header             Do not print the header
      --region string         Region for the reserved IP
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

