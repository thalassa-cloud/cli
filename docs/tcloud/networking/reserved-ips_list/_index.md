---
linkTitle: "tcloud networking reserved-ips list"
title: "networking reserved-ips list"
slug: tcloud_networking_reserved-ips_list
url: /docs/tcloud/networking/reserved-ips_list/
weight: 9792
cascade:
  type: docs
---
## tcloud networking reserved-ips list

List reserved IP addresses

### Synopsis

List reserved IP addresses within your organisation.

```
tcloud networking reserved-ips list [flags]
```

### Examples

```
tcloud networking reserved-ips list
tcloud networking reserved-ips list --region nl-ams
tcloud networking reserved-ips list --selector env=prod
```

### Options

```
      --exact-time        Show exact creation time
  -h, --help              help for list
      --no-header         Do not print the header
      --region string     Filter by region
  -l, --selector string   Label selector to filter reserved IPs (format: key1=value1,key2=value2)
      --show-labels       Show labels
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

