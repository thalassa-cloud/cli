---
linkTitle: "tcloud dns zones export"
title: "dns zones export"
slug: tcloud_dns_zones_export
url: /docs/tcloud/dns/zones_export/
weight: 9928
cascade:
  type: docs
---
## tcloud dns zones export

Export a DNS zone as a BIND zone file to stdout

```
tcloud dns zones export <zone> [flags]
```

### Options

```
  -h, --help   help for export
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

