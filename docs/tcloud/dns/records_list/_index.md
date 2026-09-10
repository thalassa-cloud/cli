---
linkTitle: "tcloud dns records list"
title: "dns records list"
slug: tcloud_dns_records_list
url: /docs/tcloud/dns/records_list/
weight: 9938
cascade:
  type: docs
---
## tcloud dns records list

List DNS records in a zone

```
tcloud dns records list [flags]
```

### Options

```
      --exact-time    Show full timestamps instead of relative time
  -h, --help          help for list
      --no-header     Do not print table headers
      --zone string   DNS zone identity
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

