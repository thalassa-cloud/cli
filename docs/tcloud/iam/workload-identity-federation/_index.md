---
linkTitle: "tcloud iam workload-identity-federation"
title: "iam workload-identity-federation"
slug: tcloud_iam_workload-identity-federation
url: /docs/tcloud/iam/workload-identity-federation/
weight: 9906
cascade:
  type: docs
---
## tcloud iam workload-identity-federation

Bootstrap and manage CI/CD workload identity (OIDC)

### Synopsis

Commands to provision federated identity providers, service accounts, role bindings,
and federated identities for GitHub Actions, GitLab CI, and Kubernetes service account OIDC tokens.

### Options

```
  -h, --help   help for workload-identity-federation
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

* [tcloud iam](/docs/tcloud/tcloud_iam/)	 - Identity and access management for your organisation
* [tcloud iam workload-identity-federation bootstrap](/docs/tcloud/iam/workload-identity-federation_bootstrap/)	 - Provision workload identity for GitHub, GitLab, or Kubernetes

