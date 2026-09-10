---
linkTitle: "tcloud networking security-groups rules set-egress"
title: "networking security-groups rules set-egress"
slug: tcloud_networking_security-groups_rules_set-egress
url: /docs/tcloud/networking/security-groups_rules_set-egress/
weight: 9772
cascade:
  type: docs
---
## tcloud networking security-groups rules set-egress

Replace all egress rules from a JSON file

```
tcloud networking security-groups rules set-egress SECURITY_GROUP [flags]
```

### Options

```
      --file string   JSON array of SecurityGroupRule objects
  -h, --help          help for set-egress
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
  -P, --project string         Project identity (overrides context; slug is resolved to identity; use "root" for organisation scope)
      --token string           Personal access token (overrides context)
```

### SEE ALSO

* [tcloud networking security-groups rules](/docs/tcloud/networking/security-groups_rules/)	 - Manage security group rules

