---
linkTitle: "tcloud quick-launch delete"
title: "quick-launch delete"
slug: tcloud_quick-launch_delete
url: /docs/tcloud/quick-launch/delete/
weight: 9710
cascade:
  type: docs
---
## tcloud quick-launch delete

Delete a quick-launch job

### Synopsis

Delete a quick-launch job.

Use --cascade Delete to also delete provisioned resources, or --cascade Orphan
to keep resources and only remove the quick-launch record.

```
tcloud quick-launch delete <identity> [flags]
```

### Options

```
      --cascade string   Child resource handling: Delete or Orphan
      --force            Skip the confirmation prompt and delete
  -h, --help             help for delete
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

