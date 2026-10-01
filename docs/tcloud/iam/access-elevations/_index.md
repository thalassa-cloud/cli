---
linkTitle: "tcloud iam access-elevations"
title: "iam access-elevations"
slug: tcloud_iam_access-elevations
url: /docs/tcloud/iam/access-elevations/
weight: 9911
cascade:
  type: docs
---
## tcloud iam access-elevations

Request and review temporary IAM access elevations

### Synopsis

Request temporary elevation to an IAM policy, and approve, reject, or revoke
requests in the current organisation/project scope.

Use --mine on list/get to work with your own requests. create and cancel always
act on your requests; approve, reject, and revoke use the approver API.

### Options

```
  -h, --help   help for access-elevations
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
* [tcloud iam access-elevations approve](/docs/tcloud/iam/access-elevations_approve/)	 - Approve a pending access elevation request
* [tcloud iam access-elevations cancel](/docs/tcloud/iam/access-elevations_cancel/)	 - Cancel your pending access elevation request
* [tcloud iam access-elevations create](/docs/tcloud/iam/access-elevations_create/)	 - Request temporary elevation to an IAM policy
* [tcloud iam access-elevations get](/docs/tcloud/iam/access-elevations_get/)	 - Show an access elevation request
* [tcloud iam access-elevations list](/docs/tcloud/iam/access-elevations_list/)	 - List access elevation requests
* [tcloud iam access-elevations reject](/docs/tcloud/iam/access-elevations_reject/)	 - Reject a pending access elevation request
* [tcloud iam access-elevations revoke](/docs/tcloud/iam/access-elevations_revoke/)	 - Revoke an approved access elevation

