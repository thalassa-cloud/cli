---
linkTitle: "tcloud networking routetables list"
title: "networking routetables list"
slug: tcloud_networking_routetables_list
url: /docs/tcloud/networking/routetables_list/
weight: 9950
cascade:
  type: docs
---
## tcloud networking routetables list

Get a list of routetables

```
tcloud networking routetables list [flags]
```

### Options

```
  -h, --help              help for list
      --no-header         Do not print the header
  -l, --selector string   Label selector to filter route tables (format: key1=value1,key2=value2)
      --show-labels       Show labels
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

* [tcloud networking routetables](/docs/tcloud/networking/routetables/)	 - Manage routetables

