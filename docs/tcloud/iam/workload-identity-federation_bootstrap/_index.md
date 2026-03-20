---
linkTitle: "tcloud iam workload-identity-federation bootstrap"
title: "iam workload-identity-federation bootstrap"
slug: tcloud_iam_workload-identity-federation_bootstrap
url: /docs/tcloud/iam/workload-identity-federation_bootstrap/
weight: 9907
cascade:
  type: docs
---
## tcloud iam workload-identity-federation bootstrap

Provision workload identity for GitHub, GitLab, or Kubernetes

### Synopsis

Creates (when missing) a federated OIDC identity provider, a Thalassa service account,
a role binding to your organisation role, and a federated identity for the workload JWT subject.

Resources are labelled thalassa.cloud/managed-by=workload-identity-bootstrap and thalassa.cloud/wif-vcs=<github|gitlab|kubernetes>.

Subcommands:
  github      — GitHub Actions (issuer https://token.actions.githubusercontent.com)
  gitlab      — GitLab CI id_token
  kubernetes  — in-cluster service account JWTs


### Options

```
      --dry-run                       Print planned changes without calling the API
  -h, --help                          help for bootstrap
      --name string                   Base name for the Thalassa service account and federated identity (federated identity becomes <name>-fi; default: wif-<platform>-<key>)
      --no-hints                      Do not print platform hints after bootstrap
      --provider-description string   Optional description when creating the federated identity provider
      --provider-name string          Optional display name when creating the federated identity provider
      --role string                   Organisation role identity, slug, or name (required)
      --scope strings                 Federated identity allowed scopes: api:read, api:write, kubernetes, objectStorage (default: api:read,api:write)
      --trusted-audience strings      JWT aud values to trust (repeatable; default: current context API URL, e.g. https://api.thalassa.cloud)
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

* [tcloud iam workload-identity-federation](/docs/tcloud/iam/workload-identity-federation/)	 - Bootstrap and manage CI/CD workload identity (OIDC)
* [tcloud iam workload-identity-federation bootstrap github](/docs/tcloud/iam/workload-identity-federation_bootstrap_github/)	 - Bootstrap workload identity for GitHub Actions
* [tcloud iam workload-identity-federation bootstrap gitlab](/docs/tcloud/iam/workload-identity-federation_bootstrap_gitlab/)	 - Bootstrap workload identity for GitLab CI
* [tcloud iam workload-identity-federation bootstrap kubernetes](/docs/tcloud/iam/workload-identity-federation_bootstrap_kubernetes/)	 - Bootstrap workload identity for Kubernetes service accounts

