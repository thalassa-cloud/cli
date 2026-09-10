---
linkTitle: "tcloud storage snapshot-policies"
title: "storage snapshot-policies"
slug: tcloud_storage_snapshot-policies
url: /docs/tcloud/storage/snapshot-policies/
weight: 9688
cascade:
  type: docs
---
## tcloud storage snapshot-policies

Manage snapshot policies

### Synopsis

Manage automated volume snapshot policies within the Thalassa Cloud Platform.

### Examples

```
tcloud storage snapshot-policies list
tcloud storage snapshot-policies create --name daily --region nl-ams --schedule '0 2 * * *' --ttl 168h --target-type selector --selector backup=true
```

### Options

```
  -h, --help   help for snapshot-policies
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

* [tcloud storage](/docs/tcloud/tcloud_storage/)	 - Manage storage resources
* [tcloud storage snapshot-policies create](/docs/tcloud/storage/snapshot-policies_create/)	 - Create a snapshot policy
* [tcloud storage snapshot-policies delete](/docs/tcloud/storage/snapshot-policies_delete/)	 - Delete snapshot policy(ies)
* [tcloud storage snapshot-policies list](/docs/tcloud/storage/snapshot-policies_list/)	 - List snapshot policies
* [tcloud storage snapshot-policies update](/docs/tcloud/storage/snapshot-policies_update/)	 - Update a snapshot policy
* [tcloud storage snapshot-policies view](/docs/tcloud/storage/snapshot-policies_view/)	 - View snapshot policy details

