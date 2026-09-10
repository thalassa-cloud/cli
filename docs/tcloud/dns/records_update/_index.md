---
linkTitle: "tcloud dns records update"
title: "dns records update"
slug: tcloud_dns_records_update
url: /docs/tcloud/dns/records_update/
weight: 9937
cascade:
  type: docs
---
## tcloud dns records update

Update a DNS record

```
tcloud dns records update <record> [flags]
```

### Options

```
      --exact-time      Show full timestamps instead of relative time
  -h, --help            help for update
      --no-header       Do not print table headers
      --ttl int         Record TTL in seconds
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

