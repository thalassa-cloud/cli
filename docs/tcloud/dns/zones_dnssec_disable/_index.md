---
linkTitle: "tcloud dns zones dnssec disable"
title: "dns zones dnssec disable"
slug: tcloud_dns_zones_dnssec_disable
url: /docs/tcloud/dns/zones_dnssec_disable/
weight: 9932
cascade:
  type: docs
---
## tcloud dns zones dnssec disable

Disable DNSSEC signing for a zone

```
tcloud dns zones dnssec disable <zone> [flags]
```

### Options

```
      --force   Skip the confirmation prompt
  -h, --help    help for disable
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

