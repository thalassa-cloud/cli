---
linkTitle: "tcloud secrets destroy-version"
title: "secrets destroy-version"
slug: tcloud_secrets_destroy-version
url: /docs/tcloud/secrets/destroy-version/
weight: 9700
cascade:
  type: docs
---
## tcloud secrets destroy-version

Permanently destroy a secret version

```
tcloud secrets destroy-version [flags]
```

### Options

```
      --force           Skip confirmation
  -h, --help            help for destroy-version
      --path string     Secret path
      --region string   Region
      --version int     Version to destroy
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

