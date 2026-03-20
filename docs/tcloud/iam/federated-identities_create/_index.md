---
linkTitle: "tcloud iam federated-identities create"
title: "iam federated-identities create"
slug: tcloud_iam_federated-identities_create
url: /docs/tcloud/iam/federated-identities_create/
weight: 9956
cascade:
  type: docs
---
## tcloud iam federated-identities create

Create a federated identity

```
tcloud iam federated-identities create [flags]
```

### Options

```
      --annotations strings          Annotations as key=value (repeatable)
      --audience-match-mode string   exact, any (default), or all
      --conditions string            Conditions as JSON object
      --conditions-file string       Path to JSON file for conditions
      --description string           Description
      --expires-at string            RFC3339 expiry time
  -h, --help                         help for create
      --labels strings               Labels as key=value (repeatable)
      --name string                  Display name
      --no-header                    Do not print table headers
      --provider string              Federated identity provider identity
      --scope strings                Allowed scopes: api:read, api:write, kubernetes, objectStorage (repeatable)
      --service-account string       Service account identity to bind
      --subject string               OIDC sub claim for this identity
      --trusted-audience strings     Trusted JWT audiences (repeatable)
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

