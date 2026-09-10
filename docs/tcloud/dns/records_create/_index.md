---
linkTitle: "tcloud dns records create"
title: "dns records create"
slug: tcloud_dns_records_create
url: /docs/tcloud/dns/records_create/
weight: 9940
cascade:
  type: docs
---
## tcloud dns records create

Create a DNS record

```
tcloud dns records create [flags]
```

### Options

```
      --exact-time      Show full timestamps instead of relative time
  -h, --help            help for create
      --name string     Record name
      --no-header       Do not print table headers
      --ttl int         Record TTL in seconds (default 300)
      --type string     Record type (TXT, A, CNAME, CAA, AAAA, MX, NS, SRV)
      --value strings   Record value (repeatable)
      --zone string     DNS zone identity
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

* [tcloud dns records](/docs/tcloud/dns/records/)	 - Manage DNS records

