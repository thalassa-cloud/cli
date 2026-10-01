---
linkTitle: "tcloud quick-launch create"
title: "quick-launch create"
slug: tcloud_quick-launch_create
url: /docs/tcloud/quick-launch/create/
weight: 9711
cascade:
  type: docs
---
## tcloud quick-launch create

Start a quick-launch job

### Synopsis

Start asynchronous provisioning from a template.

Templates:
  vpc          VPC with subnets and NAT gateway
  kubernetes  VPC stack plus a Kubernetes cluster

Provisioning continues after create returns; use get/logs to follow progress.

```
tcloud quick-launch create [flags]
```

### Examples

```
  tcloud quick-launch create --name demo --region nl-ams --template vpc
  tcloud quick-launch create --name demo-k8s --region nl-ams --template kubernetes --machine-type gp.medium
```

### Options

```
      --annotations strings   Annotations as key=value (repeatable)
      --description string    Optional description
  -h, --help                  help for create
      --labels strings        Labels as key=value (repeatable)
      --machine-type string   Node pool machine type (kubernetes template only)
      --name string           Base name prefix for created resources (required)
      --region string         Region identity, slug, or name (required)
      --subnet-cidr strings   Optional subnet CIDRs (repeatable)
      --template string       Template: vpc or kubernetes (default "vpc")
      --vpc-cidr string       Optional VPC CIDR (default assigned by the API)
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

* [tcloud quick-launch](/docs/tcloud/tcloud_quick-launch/)	 - Provision stacks from quick-launch templates

