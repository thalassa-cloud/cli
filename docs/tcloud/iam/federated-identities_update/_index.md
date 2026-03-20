---
linkTitle: "tcloud iam federated-identities update"
title: "iam federated-identities update"
slug: tcloud_iam_federated-identities_update
url: /docs/tcloud/iam/federated-identities_update/
weight: 9952
cascade:
  type: docs
---
## tcloud iam federated-identities update

Update a federated identity (only set flags are sent)

```
tcloud iam federated-identities update <identity> [flags]
```

### Options

```
      --annotations strings          Replace annotations (key=value, repeatable)
      --audience-match-mode string   exact, any, or all
      --conditions string            Conditions JSON object
      --conditions-file string       Path to JSON file for conditions
      --description string           Description
      --expires-at string            RFC3339 expiry (empty to clear not supported by all APIs)
  -h, --help                         help for update
      --labels strings               Replace labels (key=value, repeatable)
      --no-header                    Do not print table headers
      --scope strings                Replace allowed scopes (repeatable)
      --status string                active, inactive, expired, or revoked
      --trusted-audience strings     Replace trusted audiences (repeatable)
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

* [tcloud iam federated-identities](/docs/tcloud/iam/federated-identities/)	 - Federated identities (OIDC subject bindings)

