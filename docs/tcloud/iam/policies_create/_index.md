---
linkTitle: "tcloud iam policies create"
title: "iam policies create"
slug: tcloud_iam_policies_create
url: /docs/tcloud/iam/policies_create/
weight: 9895
cascade:
  type: docs
---
## tcloud iam policies create

Create an IAM policy

```
tcloud iam policies create [flags]
```

### Options

```
      --annotations strings     Annotations as key=value (repeatable)
      --description string      Policy description
  -h, --help                    help for create
      --labels strings          Labels as key=value (repeatable)
      --name string             Policy name
      --no-header               Do not print table headers
      --replicate-to-children   Replicate this policy to child projects
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

