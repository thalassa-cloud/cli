---
linkTitle: "tcloud registry namespaces configuration"
title: "registry namespaces configuration"
slug: tcloud_registry_namespaces_configuration
url: /docs/tcloud/registry/namespaces_configuration/
weight: 9718
cascade:
  type: docs
---
## tcloud registry namespaces configuration

Manage namespace configuration

### Synopsis

Manage visibility and retention policy configuration for a container registry namespace.

### Examples

```
tcloud registry namespaces configuration view --namespace crns-123
tcloud registry namespaces configuration create --namespace crns-123 --visibility private
```

### Options

```
  -h, --help   help for configuration
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

* [tcloud registry namespaces](/docs/tcloud/registry/namespaces/)	 - Manage container registry namespaces
* [tcloud registry namespaces configuration create](/docs/tcloud/registry/namespaces_configuration_create/)	 - Create namespace configuration
* [tcloud registry namespaces configuration delete](/docs/tcloud/registry/namespaces_configuration_delete/)	 - Delete namespace configuration
* [tcloud registry namespaces configuration update](/docs/tcloud/registry/namespaces_configuration_update/)	 - Update namespace configuration
* [tcloud registry namespaces configuration view](/docs/tcloud/registry/namespaces_configuration_view/)	 - View namespace configuration

