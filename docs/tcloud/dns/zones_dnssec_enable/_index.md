---
linkTitle: "tcloud dns zones dnssec enable"
title: "dns zones dnssec enable"
slug: tcloud_dns_zones_dnssec_enable
url: /docs/tcloud/dns/zones_dnssec_enable/
weight: 9931
cascade:
  type: docs
---
## tcloud dns zones dnssec enable

Enable DNSSEC signing for a zone

```
tcloud dns zones dnssec enable <zone> [flags]
```

### Options

```
  -h, --help             help for enable
      --kms-key string   KMS key identity used for DNSSEC signing
      --no-header        Do not print table headers
      --region string    Region for DNSSEC KMS key
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

