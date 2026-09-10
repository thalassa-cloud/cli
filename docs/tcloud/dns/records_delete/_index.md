---
linkTitle: "tcloud dns records delete"
title: "dns records delete"
slug: tcloud_dns_records_delete
url: /docs/tcloud/dns/records_delete/
weight: 9939
cascade:
  type: docs
---
## tcloud dns records delete

Delete a DNS record

```
tcloud dns records delete <record> [flags]
```

### Options

```
      --force         Skip the confirmation prompt and delete
  -h, --help          help for delete
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

