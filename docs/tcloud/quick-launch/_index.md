---
linkTitle: "tcloud quick-launch"
title: "quick-launch"
slug: tcloud_quick-launch
url: /docs/tcloud/tcloud_quick-launch/
weight: 9706
cascade:
  type: docs
---
## tcloud quick-launch

Provision stacks from quick-launch templates

### Synopsis

Start and manage quick-launch jobs that provision common stacks asynchronously
(VPC networking, or VPC + Kubernetes) in the current organisation/project scope.

### Options

```
  -h, --help   help for quick-launch
```

### Options inherited from parent commands

```
      --access-token string    Access Token authentication (overrides context)
      --api string             API endpoint (overrides context)
      --client-id string       OIDC client ID for OIDC authentication (overrides context)
      --client-secret string   OIDC client secret for OIDC authentication (overrides context)
  -c, --context string         Context name
      --debug                  Debug mode
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud quick-launch create](/docs/tcloud/quick-launch/create/)	 - Start a quick-launch job
* [tcloud quick-launch delete](/docs/tcloud/quick-launch/delete/)	 - Delete a quick-launch job
* [tcloud quick-launch get](/docs/tcloud/quick-launch/get/)	 - Show a quick-launch job
* [tcloud quick-launch list](/docs/tcloud/quick-launch/list/)	 - List quick-launch jobs
* [tcloud quick-launch logs](/docs/tcloud/quick-launch/logs/)	 - Show logs for a quick-launch job

