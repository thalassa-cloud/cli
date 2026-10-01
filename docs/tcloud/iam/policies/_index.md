---
linkTitle: "tcloud iam policies"
title: "iam policies"
slug: tcloud_iam_policies
url: /docs/tcloud/iam/policies/
weight: 9887
cascade:
  type: docs
---
## tcloud iam policies

IAM policies, permission rules, and bindings

### Synopsis

IAM policies define permission rules and can be bound to users or service accounts
in the current organisation/project scope. System policies may be read-only; the API
enforces what you can change.

### Options

```
  -h, --help   help for policies
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

* [tcloud iam](/docs/tcloud/tcloud_iam/)	 - Identity and access management for your organisation
* [tcloud iam policies bindings](/docs/tcloud/iam/policies_bindings/)	 - Policy bindings (who receives the policy)
* [tcloud iam policies create](/docs/tcloud/iam/policies_create/)	 - Create an IAM policy
* [tcloud iam policies delete](/docs/tcloud/iam/policies_delete/)	 - Delete an IAM policy
* [tcloud iam policies get](/docs/tcloud/iam/policies_get/)	 - Show a policy including rules and bindings summary
* [tcloud iam policies list](/docs/tcloud/iam/policies_list/)	 - List IAM policies
* [tcloud iam policies rules](/docs/tcloud/iam/policies_rules/)	 - Permission rules on a policy
* [tcloud iam policies update](/docs/tcloud/iam/policies_update/)	 - Update an IAM policy (only set flags are changed)

