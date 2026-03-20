---
linkTitle: "tcloud iam federated-identity-providers create"
title: "iam federated-identity-providers create"
slug: tcloud_iam_federated-identity-providers_create
url: /docs/tcloud/iam/federated-identity-providers_create/
weight: 9950
cascade:
  type: docs
---
## tcloud iam federated-identity-providers create

Register a federated identity provider

```
tcloud iam federated-identity-providers create [flags]
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Description
  -h, --help                  help for create
      --issuer string         OIDC issuer URL (unique per organisation)
      --jwks-uri string       Optional JWKS URI override
      --labels strings        Labels as key=value (repeatable)
      --name string           Provider name
      --no-header             Do not print table headers
      --status string         active (default) or inactive
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam federated-identity-providers](/docs/tcloud/iam/federated-identity-providers/)	 - Federated OIDC identity providers

