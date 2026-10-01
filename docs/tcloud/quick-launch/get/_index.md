---
linkTitle: "tcloud quick-launch get"
title: "quick-launch get"
slug: tcloud_quick-launch_get
url: /docs/tcloud/quick-launch/get/
weight: 9709
cascade:
  type: docs
---
## tcloud quick-launch get

Show a quick-launch job

```
tcloud quick-launch get <identity> [flags]
```

### Options

```
      --exact-time   Show full timestamps instead of relative time
  -h, --help         help for get
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

* [tcloud quick-launch](/docs/tcloud/tcloud_quick-launch/)	 - Provision stacks from quick-launch templates

