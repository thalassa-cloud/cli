---
linkTitle: "tcloud iam federated-identity-providers update"
title: "iam federated-identity-providers update"
slug: tcloud_iam_federated-identity-providers_update
url: /docs/tcloud/iam/federated-identity-providers_update/
weight: 9946
cascade:
  type: docs
---
## tcloud iam federated-identity-providers update

Update a federated identity provider

```
tcloud iam federated-identity-providers update <identity> [flags]
```

### Options

```
      --annotations strings   Replace annotations (key=value, repeatable)
      --description string    Description
  -h, --help                  help for update
      --jwks-uri string       JWKS URI (set to empty string to clear)
      --labels strings        Replace labels (key=value, repeatable)
      --name string           Provider display name
      --no-header             Do not print table headers
      --status string         active or inactive
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

