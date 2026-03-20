---
linkTitle: "tcloud iam members delete"
title: "iam members delete"
slug: tcloud_iam_members_delete
url: /docs/tcloud/iam/members_delete/
weight: 9942
cascade:
  type: docs
---
## tcloud iam members delete

Remove a member from the organisation

```
tcloud iam members delete <member> [flags]
```

### Options

```
      --force       Skip the confirmation prompt and remove the member
  -h, --help        help for delete
      --no-header   Do not print table headers
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

* [tcloud iam members](/docs/tcloud/iam/members/)	 - Organisation members (owners and members)

