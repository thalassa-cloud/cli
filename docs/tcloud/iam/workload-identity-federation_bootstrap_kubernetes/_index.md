---
linkTitle: "tcloud iam workload-identity-federation bootstrap kubernetes"
title: "iam workload-identity-federation bootstrap kubernetes"
slug: tcloud_iam_workload-identity-federation_bootstrap_kubernetes
url: /docs/tcloud/iam/workload-identity-federation_bootstrap_kubernetes/
weight: 9908
cascade:
  type: docs
---
## tcloud iam workload-identity-federation bootstrap kubernetes

Bootstrap workload identity for Kubernetes service accounts

### Synopsis

Binds system:serviceaccount:<namespace>:<name> to a Thalassa service account via a federated identity.

Thalassa clusters: pass --cluster to resolve the cluster and use the platform-managed federated identity
provider labelled kubernetes_cluster_id=<cluster identity>. If that provider is missing, bootstrap fails
(you must wait for cluster OIDC integration).

Other clusters: pass --issuer with the same URL as kube-apiserver --service-account-issuer; the CLI
creates the federated identity provider if it does not exist yet.

```
tcloud iam workload-identity-federation bootstrap kubernetes [flags]
```

### Examples

```
  # Thalassa-managed cluster (OIDC provider from label kubernetes_cluster_id)
  tcloud iam workload-identity-federation bootstrap kubernetes --cluster my-cluster-slug \
    --namespace default --service-account my-app --role deployer

  # Self-managed / custom issuer
  tcloud iam workload-identity-federation bootstrap kubernetes --issuer https://k8s.example.com \
    --namespace cicd --service-account terraform --role deployer
```

### Options

```
      --cluster string           Thalassa cluster identity, slug, or name (uses federated identity provider with label kubernetes_cluster_id=<cluster identity>)
  -h, --help                     help for kubernetes
      --issuer string            kube-apiserver service-account issuer URL (required without --cluster); with --cluster must match the cluster OIDC provider or omit
      --namespace string         Kubernetes namespace for the workload service account (required)
      --service-account string   Kubernetes service account name (required)
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

