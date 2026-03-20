---
linkTitle: "tcloud iam teams update"
title: "iam teams update"
slug: tcloud_iam_teams_update
url: /docs/tcloud/iam/teams_update/
weight: 9912
cascade:
  type: docs
---
## tcloud iam teams update

Update a team (only flags you set are changed)

```
tcloud iam teams update <team> [flags]
```

### Options

```
      --annotations strings   Replace annotations (key=value, repeatable)
      --description string    Team description
      --exact-time            Show full timestamps instead of relative time
  -h, --help                  help for update
      --labels strings        Replace labels (key=value, repeatable)
      --name string           Team display name
      --no-header             Do not print table headers
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

* [tcloud iam teams](/docs/tcloud/iam/teams/)	 - Manage organisation teams

