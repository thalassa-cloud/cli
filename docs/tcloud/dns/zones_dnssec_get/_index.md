---
linkTitle: "tcloud dns zones dnssec get"
title: "dns zones dnssec get"
slug: tcloud_dns_zones_dnssec_get
url: /docs/tcloud/dns/zones_dnssec_get/
weight: 9930
cascade:
  type: docs
---
## tcloud dns zones dnssec get

Get DNSSEC status for a zone

```
tcloud dns zones dnssec get <zone> [flags]
```

### Options

```
  -h, --help        help for get
      --no-header   Do not print table headers
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

* [tcloud dns zones dnssec](/docs/tcloud/dns/zones_dnssec/)	 - Manage DNSSEC for a DNS zone

