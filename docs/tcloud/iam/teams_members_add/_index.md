---
linkTitle: "tcloud iam teams members add"
title: "iam teams members add"
slug: tcloud_iam_teams_members_add
url: /docs/tcloud/iam/teams_members_add/
weight: 9916
cascade:
  type: docs
---
## tcloud iam teams members add

Add a user to a team

```
tcloud iam teams members add <team> [flags]
```

### Options

```
  -h, --help          help for add
      --no-header     Do not print table headers
      --role string   Team role for the user
      --user string   User identity to add
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

* [tcloud iam teams members](/docs/tcloud/iam/teams_members/)	 - Manage team membership

