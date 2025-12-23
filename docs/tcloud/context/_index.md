---
linkTitle: "tcloud context"
title: "context"
slug: tcloud_context
url: /docs/tcloud/tcloud_context/
weight: 9979
cascade:
  type: docs
---
## tcloud context

Manage context

### Synopsis

Manage context for the CLI. Contexts are used to manage multiple organisations and APIs.

### Options

```
  -h, --help   help for context
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

* [tcloud](/docs/tcloud/tcloud/)	 - A CLI for working with the Thalassa Cloud Platform
* [tcloud context create](/docs/tcloud/context/create/)	 - Create a new context with authentication and organisation
* [tcloud context current](/docs/tcloud/context/current/)	 - Shows the current context
* [tcloud context delete](/docs/tcloud/context/delete/)	 - Delete a context
* [tcloud context delete-server](/docs/tcloud/context/delete-server/)	 - Delete a server
* [tcloud context delete-user](/docs/tcloud/context/delete-user/)	 - Delete a user
* [tcloud context list](/docs/tcloud/context/list/)	 - List the contexts
* [tcloud context login](/docs/tcloud/context/login/)	 - Login to Thalassa Cloud
* [tcloud context organisation](/docs/tcloud/context/organisation/)	 - Set the organisation in the current-context
* [tcloud context use](/docs/tcloud/context/use/)	 - Set the current context
* [tcloud context view](/docs/tcloud/context/view/)	 - Shows current context

