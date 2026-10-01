---
linkTitle: "tcloud quick-launch logs"
title: "quick-launch logs"
slug: tcloud_quick-launch_logs
url: /docs/tcloud/quick-launch/logs/
weight: 9707
cascade:
  type: docs
---
## tcloud quick-launch logs

Show logs for a quick-launch job

```
tcloud quick-launch logs <identity> [flags]
```

### Options

```
  -h, --help   help for logs
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

