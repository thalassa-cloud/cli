---
linkTitle: "tcloud context fix"
title: "context fix"
slug: tcloud_context_fix
url: /docs/tcloud/context/fix/
weight: 9982
cascade:
  type: docs
---
## tcloud context fix

Fix config file security issues

### Synopsis

Fix security issues in the CLI config file, such as overly permissive file permissions or migrating credentials to the keychain.

```
tcloud context fix [flags]
```

### Examples

```
  tcloud context fix
  tcloud context fix --migrate-credentials
```

### Options

```
  -h, --help                  help for fix
      --migrate-credentials   move plaintext credentials from the config file into the keychain
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

* [tcloud context](/docs/tcloud/tcloud_context/)	 - Manage context

