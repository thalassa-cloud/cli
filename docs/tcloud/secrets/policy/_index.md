---
linkTitle: "tcloud secrets policy"
title: "secrets policy"
slug: tcloud_secrets_policy
url: /docs/tcloud/secrets/policy/
weight: 9697
cascade:
  type: docs
---
## tcloud secrets policy

Replace a secret access policy from a JSON file

```
tcloud secrets policy [flags]
```

### Options

```
      --file string     JSON file containing the access policy
  -h, --help            help for policy
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

