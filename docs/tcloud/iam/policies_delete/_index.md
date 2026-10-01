---
linkTitle: "tcloud iam policies delete"
title: "iam policies delete"
slug: tcloud_iam_policies_delete
url: /docs/tcloud/iam/policies_delete/
weight: 9894
cascade:
  type: docs
---
## tcloud iam policies delete

Delete an IAM policy

```
tcloud iam policies delete <policy> [flags]
```

### Options

```
      --force       Skip the confirmation prompt and delete
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
      --ignore-dir-config      Ignore directory-local .thalassa defaults
  -O, --organisation string    Organisation slug or identity (overrides context)
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud iam policies](/docs/tcloud/iam/policies/)	 - IAM policies, permission rules, and bindings

