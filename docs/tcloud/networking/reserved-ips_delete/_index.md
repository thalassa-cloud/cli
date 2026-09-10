---
linkTitle: "tcloud networking reserved-ips delete"
title: "networking reserved-ips delete"
slug: tcloud_networking_reserved-ips_delete
url: /docs/tcloud/networking/reserved-ips_delete/
weight: 9794
cascade:
  type: docs
---
## tcloud networking reserved-ips delete

Delete reserved IP address(es)

### Synopsis

Delete reserved IP address(es) by identity or label selector. Attached reserved IPs are disassociated before deletion.

```
tcloud networking reserved-ips delete [flags]
```

### Examples

```
tcloud networking reserved-ips delete rip-123 --force
tcloud networking reserved-ips delete --selector env=test --force
```

### Options

```
      --force             Force the deletion and skip the confirmation
  -h, --help              help for delete
  -l, --selector string   Label selector to filter reserved IPs (format: key1=value1,key2=value2)
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

