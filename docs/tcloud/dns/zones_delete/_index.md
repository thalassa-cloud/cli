---
linkTitle: "tcloud dns zones delete"
title: "dns zones delete"
slug: tcloud_dns_zones_delete
url: /docs/tcloud/dns/zones_delete/
weight: 9933
cascade:
  type: docs
---
## tcloud dns zones delete

Delete a DNS zone and all of its records

```
tcloud dns zones delete <zone> [flags]
```

### Options

```
      --force   Skip the confirmation prompt and delete
  -h, --help    help for delete
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

* [tcloud dns zones](/docs/tcloud/dns/zones/)	 - Manage DNS zones

