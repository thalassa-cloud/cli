---
linkTitle: "tcloud iam access-elevations create"
title: "iam access-elevations create"
slug: tcloud_iam_access-elevations_create
url: /docs/tcloud/iam/access-elevations_create/
weight: 9916
cascade:
  type: docs
---
## tcloud iam access-elevations create

Request temporary elevation to an IAM policy

### Synopsis

Create an access elevation request for yourself against a policy.

Exactly one of --duration or --expires-at is required.

```
tcloud iam access-elevations create [flags]
```

### Examples

```
  tcloud iam access-elevations create --policy admin --reason "incident response" --duration 2h
  tcloud iam access-elevations create --policy pol-xxx --reason "break-glass" --expires-at 2026-10-01T20:00:00Z
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --duration string       Relative TTL (e.g. 2h, 24h)
      --expires-at string     Absolute expiry time (RFC3339)
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --policy string         IAM policy identity, slug, or name (required)
      --reason string         Reason for the elevation request (required)
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

* [tcloud iam access-elevations](/docs/tcloud/iam/access-elevations/)	 - Request and review temporary IAM access elevations

