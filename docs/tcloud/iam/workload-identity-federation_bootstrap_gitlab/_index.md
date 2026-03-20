---
linkTitle: "tcloud iam workload-identity-federation bootstrap gitlab"
title: "iam workload-identity-federation bootstrap gitlab"
slug: tcloud_iam_workload-identity-federation_bootstrap_gitlab
url: /docs/tcloud/iam/workload-identity-federation_bootstrap_gitlab/
weight: 9909
cascade:
  type: docs
---
## tcloud iam workload-identity-federation bootstrap gitlab

Bootstrap workload identity for GitLab CI

### Synopsis

Creates or reuses a federated identity provider for your GitLab OIDC issuer, then binds the project
and ref to a Thalassa service account.

The GitLab id_token sub uses project_path:<group/project>:ref_type:<type>:ref:<ref>.

```
tcloud iam workload-identity-federation bootstrap gitlab [flags]
```

### Examples

```
  # GitLab.com, branch main
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref main --role deployer

  # Tag pipeline
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref v1.0.0 --ref-type tag --role deployer

  # Self-managed GitLab
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref main --issuer https://gitlab.example.com --role deployer
```

### Options

```
  -h, --help                help for gitlab
      --issuer string       GitLab OIDC issuer URL (self-managed) (default "https://gitlab.com")
      --ref string          Git ref segment for id_token sub, e.g. branch or tag name (required)
      --ref-type string     ref_type claim in JWT sub: branch, tag, merge_request, etc. (default "branch")
      --repository string   GitLab project path as group/project (required)
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

