---
linkTitle: "tcloud dns zones create"
title: "dns zones create"
slug: tcloud_dns_zones_create
url: /docs/tcloud/dns/zones_create/
weight: 9934
cascade:
  type: docs
---
## tcloud dns zones create

Create a DNS zone

```
tcloud dns zones create [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Zone description
      --exact-time            Show full timestamps instead of relative time
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --name string           Zone name (e.g. example.com)
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

