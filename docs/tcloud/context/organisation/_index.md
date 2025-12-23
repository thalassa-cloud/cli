---
linkTitle: "tcloud context organisation"
title: "context organisation"
slug: tcloud_context_organisation
url: /docs/tcloud/context/organisation/
weight: 9982
cascade:
  type: docs
---
## tcloud context organisation

Set the organisation in the current-context

### Synopsis

Set the organisation in the current-context

```
tcloud context organisation <organisation> [flags]
```

### Examples

```
tcloud context use-organisation <organisation>
```

### Options

```
  -h, --help   help for organisation
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
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

