---
linkTitle: "tcloud dns zones update"
title: "dns zones update"
slug: tcloud_dns_zones_update
url: /docs/tcloud/dns/zones_update/
weight: 9925
cascade:
  type: docs
---
## tcloud dns zones update

Update a DNS zone

```
tcloud dns zones update <zone> [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Zone description
      --exact-time            Show full timestamps instead of relative time
  -h, --help                  help for update
      --labels strings        Labels as key=value (repeatable)
      --no-header             Do not print table headers
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

