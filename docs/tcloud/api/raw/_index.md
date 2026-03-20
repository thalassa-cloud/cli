---
linkTitle: "tcloud api raw"
title: "api raw"
slug: tcloud_api_raw
url: /docs/tcloud/api/raw/
weight: 9999
cascade:
  type: docs
---
## tcloud api raw

Make a raw HTTP request to the API

### Synopsis

Make a raw HTTP request to the Thalassa Cloud API.

Similar to 'kubectl get --raw', this bypasses the CLI resource layer and sends
the request directly to the API server. Uses the same authentication and
context (organisation, endpoint) as other tcloud commands.

PATH must start with a slash (e.g. /v1/me/organisations).
Requires client-go with RawRequest support.

```
tcloud api raw PATH [flags]
```

### Examples

```
  tcloud api raw /v1/me/organisations
  tcloud api raw -X GET /v1/iaas/regions
  tcloud api raw -X POST -d '{"name":"test"}' /v1/some/resource
  tcloud api raw --show-headers /v1/me
```

### Options

```
  -d, --data string      Request body (for POST, PUT, PATCH)
  -h, --help             help for raw
  -X, --request string   HTTP method (GET, POST, PUT, PATCH, DELETE) (default "GET")
      --show-headers     Print response headers
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

* [tcloud api](/docs/tcloud/tcloud_api/)	 - Direct API access

