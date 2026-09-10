---
linkTitle: "tcloud secrets browse"
title: "secrets browse"
slug: tcloud_secrets_browse
url: /docs/tcloud/secrets/browse/
weight: 9703
cascade:
  type: docs
---
## tcloud secrets browse

Browse secret prefixes and secrets at a path

### Synopsis

Browse Secrets Manager prefixes and secrets.

In a terminal with fzf available, browse is interactive: select prefixes to
descend, ".." to go up, and secrets to view metadata (optionally reveal values
after confirmation). Press Esc to quit.

When stdout is not a terminal, or TC_IGNORE_FZF is set, prints a non-interactive
table for the given --path.

```
tcloud secrets browse [flags]
```

### Options

```
      --exact-time      Show full timestamps instead of relative time
  -h, --help            help for browse
      --no-header       Do not print table headers
      --path string     Path to browse (default "/")
      --region string   Region
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

* [tcloud secrets](/docs/tcloud/tcloud_secrets/)	 - Manage secrets (beta)

