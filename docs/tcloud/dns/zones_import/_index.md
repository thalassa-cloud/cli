---
linkTitle: "tcloud dns zones import"
title: "dns zones import"
slug: tcloud_dns_zones_import
url: /docs/tcloud/dns/zones_import/
weight: 9927
cascade:
  type: docs
---
## tcloud dns zones import

Import DNS records from a BIND zone file

```
tcloud dns zones import <zone> [flags]
```

### Options

```
      --file string   Path to BIND zone file
  -h, --help          help for import
      --no-header     Do not print table headers
      --replace       Replace existing records that conflict
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

