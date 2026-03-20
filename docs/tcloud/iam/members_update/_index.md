---
linkTitle: "tcloud iam members update"
title: "iam members update"
slug: tcloud_iam_members_update
url: /docs/tcloud/iam/members_update/
weight: 9940
cascade:
  type: docs
---
## tcloud iam members update

Change an organisation member's role (OWNER or MEMBER)

```
tcloud iam members update <member> [flags]
```

### Options

```
  -h, --help          help for update
      --no-header     Do not print table headers
      --role string   Organisation role: OWNER or MEMBER
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

