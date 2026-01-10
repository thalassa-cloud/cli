---
linkTitle: "tcloud oidc token-exchange"
title: "oidc token-exchange"
slug: tcloud_oidc_token-exchange
url: /docs/tcloud/oidc/token-exchange/
weight: 9922
cascade:
  type: docs
---
## tcloud oidc token-exchange

Exchange an OIDC token for an access token

### Synopsis

Helper for exchanging an OIDC subject token for a Thalassa Cloud access token using the token exchange grant type. Intended for use in CI/CD pipelines such as GitLab CI, GitHub Actions, Kubernetes, and similar automation environments where OIDC identity tokens are provided and need to be exchanged for a Thalassa access token.

```
tcloud oidc token-exchange [flags]
```

### Options

```
      --access-token-lifetime string   Access token lifetime (min: 1m, max: 24h, default: 1h) (default "1h")
  -h, --help                           help for token-exchange
      --organisation-id string         Organisation ID (can also be set via context)
      --service-account-id string      Service account ID (can also be set via THALASSA_SERVICE_ACCOUNT_ID env var)
      --subject-token string           Subject token (JWT) to exchange (can also be set via THALASSA_ID_TOKEN env var)
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

* [tcloud oidc](/docs/tcloud/tcloud_oidc/)	 - OIDC token operations

