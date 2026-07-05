---
linkTitle: "tcloud registry namespaces configuration create"
title: "registry namespaces configuration create"
slug: tcloud_registry_namespaces_configuration_create
url: /docs/tcloud/registry/namespaces_configuration_create/
weight: 9808
cascade:
  type: docs
---
## tcloud registry namespaces configuration create

Create namespace configuration

```
tcloud registry namespaces configuration create [flags]
```

### Examples

```
tcloud registry namespaces configuration create --namespace crns-123 --visibility private --retention-enabled --retention-days 30
```

### Options

```
      --delete-untagged                Delete untagged images during retention runs
  -h, --help                           help for create
      --namespace string               Namespace identity
      --retention-count int            Retain this many recent tags (simple single-rule policy)
      --retention-days int             Retain tags for this many days (simple single-rule policy)
      --retention-enabled              Enable retention policy
      --retention-policy-file string   Path to JSON file with full retention policy
      --visibility string              Namespace visibility (default "private")
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud registry namespaces configuration](/docs/tcloud/registry/namespaces_configuration/)	 - Manage namespace configuration

