---
linkTitle: "tcloud dns records view"
title: "dns records view"
slug: tcloud_dns_records_view
url: /docs/tcloud/dns/records_view/
weight: 9936
cascade:
  type: docs
---
## tcloud dns records view

View a DNS record

```
tcloud dns records view <record> [flags]
```

### Options

```
      --exact-time    Show full timestamps instead of relative time
  -h, --help          help for view
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

