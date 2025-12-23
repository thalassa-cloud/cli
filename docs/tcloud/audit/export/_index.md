---
linkTitle: "tcloud audit export"
title: "audit export"
slug: tcloud_audit_export
url: /docs/tcloud/audit/export/
weight: 9999
cascade:
  type: docs
---
## tcloud audit export

Export organisation audit logs to a JSON file

### Synopsis

Export organisation audit logs to a JSON file for compliance purposes.

Time Range Options:
  Use --since to specify a duration (e.g., --since 364d, --since 1w, --since 1mo, --since 1y)
  Or use --from/--to for explicit date ranges

Split Options:
  Use --daily, --weekly, or --monthly to split the export into separate files per period
  When using splits, each period is exported to a separate file

Output:
  Use '-' as the output file to write to stdout

Examples:
  tcloud audit export --since 7d --daily
  tcloud audit export --since 364d --weekly
  tcloud audit export --since 30d --monthly
  tcloud audit export --from 2024-01-01 --to 2024-01-31 --daily
  tcloud audit export --since 1d --output audit-logs.json
  tcloud audit export --since 1d --output -

```
tcloud audit export [flags]
```

### Options

```
      --action strings                  Filter by action(s) (can be specified multiple times)
      --chunk-download-timeout string   API call timeout per chunk (e.g., 5m, 10m, 1h, default: 5m)
      --daily                           Split export into separate files per day
      --from string                     Start date for custom range (YYYY-MM-DD)
  -h, --help                            help for export
      --impersonator-identity string    Filter by impersonator identity
      --include-system-services         Include system service logs
      --monthly                         Split export into separate files per month
      --organisation-identity string    Filter by organisation identity
  -o, --output string                   Output file path (use '-' for stdout, default: audit-logs-{range}-{timestamp}.json)
      --page-size int                   Page size for audit log list (default: 100) (default 100)
      --resource-identity string        Filter by resource identity
      --resource-type strings           Filter by resource type(s) (can be specified multiple times)
      --response-status int             Filter by HTTP response status code
      --search-text string              Search text filter
      --service-account string          Filter by service account identity
      --since string                    Export logs from the past duration (e.g., 364d, 1w, 1mo, 1y, 24h)
      --to string                       End date for custom range (YYYY-MM-DD)
      --user-identity string            Filter by user identity
      --weekly                          Split export into separate files per week
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

* [tcloud audit](/docs/tcloud/tcloud_audit/)	 - Manage organisation audit logs

