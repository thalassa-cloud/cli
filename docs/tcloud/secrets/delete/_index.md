---
linkTitle: "tcloud secrets delete"
title: "secrets delete"
slug: tcloud_secrets_delete
url: /docs/tcloud/secrets/delete/
weight: 9701
cascade:
  type: docs
---
## tcloud secrets delete

Delete a secret

```
tcloud secrets delete [flags]
```

### Options

```
      --force           Skip confirmation
  -h, --help            help for delete
      --path string     Secret path
      --region string   Region
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

* [tcloud secrets](/docs/tcloud/tcloud_secrets/)	 - Manage secrets (beta)

