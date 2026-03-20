---
linkTitle: "tcloud iam"
title: "iam"
slug: tcloud_iam
url: /docs/tcloud/tcloud_iam/
weight: 9905
cascade:
  type: docs
---
## tcloud iam

Identity and access management for your organisation

### Synopsis

Manage teams, organisation members, custom roles, federated OIDC identities,
and related resources. Commands apply to the organisation selected in your context
(or the --organisation / -O flag).

### Options

```
  -h, --help   help for iam
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud iam federated-identities](/docs/tcloud/iam/federated-identities/)	 - Federated identities (OIDC subject bindings)
* [tcloud iam federated-identity-providers](/docs/tcloud/iam/federated-identity-providers/)	 - Federated OIDC identity providers
* [tcloud iam invites](/docs/tcloud/iam/invites/)	 - Organisation member invitations
* [tcloud iam members](/docs/tcloud/iam/members/)	 - Organisation members (owners and members)
* [tcloud iam roles](/docs/tcloud/iam/roles/)	 - Organisation IAM roles, permission rules, and bindings
* [tcloud iam service-accounts](/docs/tcloud/iam/service-accounts/)	 - Organisation service accounts
* [tcloud iam teams](/docs/tcloud/iam/teams/)	 - Manage organisation teams
* [tcloud iam workload-identity-federation](/docs/tcloud/iam/workload-identity-federation/)	 - Bootstrap and manage CI/CD workload identity (OIDC)

