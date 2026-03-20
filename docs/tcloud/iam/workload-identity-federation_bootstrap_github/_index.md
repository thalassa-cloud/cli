---
linkTitle: "tcloud iam workload-identity-federation bootstrap github"
title: "iam workload-identity-federation bootstrap github"
slug: tcloud_iam_workload-identity-federation_bootstrap_github
url: /docs/tcloud/iam/workload-identity-federation_bootstrap_github/
weight: 9910
cascade:
  type: docs
---
## tcloud iam workload-identity-federation bootstrap github

Bootstrap workload identity for GitHub Actions

### Synopsis

Creates or reuses a federated identity provider for GitHub's OIDC issuer, then binds your repository
and ref to a Thalassa service account via a federated identity.

The JWT issuer is https://token.actions.githubusercontent.com. Match subjects with --repository (owner/name),
--ref-kind, and --ref (omit --ref only when --ref-kind is pull_request).

```
tcloud iam workload-identity-federation bootstrap github [flags]
```

### Examples

```
  # Main branch (JWT aud defaults to context API URL)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref main --role deployer

  # Specific ref kind
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref-kind branch --ref main --role deployer

  # Pull request workflows (subject repo:owner/repo:pull_request)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref-kind pull_request --role deployer

  # Custom Thalassa service account / federated identity names (federated identity: platform-ci-fi)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref main --name platform-ci --role deployer
```

### Options

```
  -h, --help                help for github
      --ref string          Branch, tag, or environment name (omit when --ref-kind pull_request)
      --ref-kind string     branch, tag, environment, or pull_request (default "branch")
      --repository string   GitHub repository as owner/name (required)
```

### Options inherited from parent commands

```
      --access-token string           Access Token authentication (overrides context)
      --api string                    API endpoint (overrides context)
      --client-id string              OIDC client ID for OIDC authentication (overrides context)
      --client-secret string          OIDC client secret for OIDC authentication (overrides context)
  -c, --context string                Context name
      --debug                         Debug mode
      --dry-run                       Print planned changes without calling the API
      --name string                   Base name for the Thalassa service account and federated identity (federated identity becomes <name>-fi; default: wif-<platform>-<key>)
      --no-hints                      Do not print platform hints after bootstrap
  -O, --organisation string           Organisation slug or identity (overrides context)
      --provider-description string   Optional description when creating the federated identity provider
      --provider-name string          Optional display name when creating the federated identity provider
      --role string                   Organisation role identity, slug, or name (required)
      --scope strings                 Federated identity allowed scopes: api:read, api:write, kubernetes, objectStorage (default: api:read,api:write)
      --token string                  Personal access token (overrides context)
      --trusted-audience strings      JWT aud values to trust (repeatable; default: current context API URL, e.g. https://api.thalassa.cloud)
```

### SEE ALSO

* [tcloud iam workload-identity-federation bootstrap](/docs/tcloud/iam/workload-identity-federation_bootstrap/)	 - Provision workload identity for GitHub, GitLab, or Kubernetes

