---
linkTitle: "tcloud registry namespaces"
title: "registry namespaces"
slug: tcloud_registry_namespaces
url: /docs/tcloud/registry/namespaces/
weight: 9710
cascade:
  type: docs
---
## tcloud registry namespaces

Manage container registry namespaces

### Examples

```
tcloud registry namespaces list
tcloud registry namespaces create --namespace my-app --region eu-west-1
```

### Options

```
  -h, --help   help for namespaces
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

* [tcloud registry](/docs/tcloud/tcloud_registry/)	 - Manage the Thalassa container registry
* [tcloud registry namespaces configuration](/docs/tcloud/registry/namespaces_configuration/)	 - Manage namespace configuration
* [tcloud registry namespaces create](/docs/tcloud/registry/namespaces_create/)	 - Create a container registry namespace
* [tcloud registry namespaces delete](/docs/tcloud/registry/namespaces_delete/)	 - Delete container registry namespace(s)
* [tcloud registry namespaces list](/docs/tcloud/registry/namespaces_list/)	 - List container registry namespaces
* [tcloud registry namespaces retention](/docs/tcloud/registry/namespaces_retention/)	 - Run retention policies
* [tcloud registry namespaces update](/docs/tcloud/registry/namespaces_update/)	 - Update a container registry namespace
* [tcloud registry namespaces view](/docs/tcloud/registry/namespaces_view/)	 - View namespace details

