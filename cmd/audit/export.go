package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/audit"
	"github.com/thalassa-cloud/client-go/thalassa"
)

var (
	outputFile            string
	daily                 bool
	weekly                bool
	monthly               bool
	splitPeriod           string
	sinceDuration         string
	startDate             string
	endDate               string
	apiTimeout            string
	searchText            string
	serviceAccount        string
	userIdentity          string
	impersonatorIdentity  string
	actions               []string
	resourceTypes         []string
	resourceIdentity      string
	organisationIdentity  string
	includeSystemServices bool
	responseStatus        int

	pageSize int
)

// parseDuration parses a duration string supporting days (d), weeks (w), months (mo), years (y)
// in addition to standard Go duration units (h, m, s, etc.)
func parseDuration(s string) (time.Duration, error) {
	// First try standard time.ParseDuration
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Handle custom units: d, w, mo, y
	var total time.Duration
	re := regexp.MustCompile(`(\d+)([dwy]|mo)`)
	matches := re.FindAllStringSubmatch(s, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", match[1])
		}

		unit := match[2]
		switch unit {
		case "d":
			total += time.Duration(value) * 24 * time.Hour
		case "w":
			total += time.Duration(value) * 7 * 24 * time.Hour
		case "mo":
			// Approximate month as 30 days
			total += time.Duration(value) * 30 * 24 * time.Hour
		case "y":
			// Approximate year as 365 days
			total += time.Duration(value) * 365 * 24 * time.Hour
		}
	}

	// Check if there are remaining characters not matched
	remaining := re.ReplaceAllString(s, "")
	if strings.TrimSpace(remaining) != "" {
		return 0, fmt.Errorf("invalid duration format: %s (unrecognized: %s)", s, remaining)
	}

	return total, nil
}

// splitTimeRange splits a time range into chunks based on the split period
func splitTimeRange(start, end time.Time, splitPeriod string) [][]time.Time {
	var chunks [][]time.Time
	current := start

	for current.Before(end) {
		chunkStart := current
		var chunkEnd time.Time

		switch splitPeriod {
		case "daily":
			chunkEnd = time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 999999999, current.Location())
			current = chunkEnd.Add(time.Second)
			// Move to start of next day
			current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		case "weekly":
			// Find end of week (Sunday)
			weekday := int(current.Weekday())
			if weekday == 0 {
				weekday = 7 // Sunday is 7
			}
			daysUntilSunday := 7 - weekday
			chunkEnd = current.AddDate(0, 0, daysUntilSunday)
			chunkEnd = time.Date(chunkEnd.Year(), chunkEnd.Month(), chunkEnd.Day(), 23, 59, 59, 999999999, chunkEnd.Location())
			// Move to start of next Monday
			current = chunkEnd.AddDate(0, 0, 1)
			current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		case "monthly":
			// End of current month
			chunkEnd = time.Date(current.Year(), current.Month()+1, 0, 23, 59, 59, 999999999, current.Location())
			current = chunkEnd.Add(time.Second)
			// Move to start of next month
			current = time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
		default:
			// No split, return single chunk
			return [][]time.Time{{start, end}}
		}

		// Don't exceed the end time
		if chunkEnd.After(end) {
			chunkEnd = end
		}

		chunks = append(chunks, []time.Time{chunkStart, chunkEnd})

		// Stop if we've reached the end
		if !chunkEnd.Before(end) || current.After(end) {
			break
		}
	}

	return chunks
}

// exportTimeRange exports audit logs for a specific time range
func exportTimeRange(ctx context.Context, client thalassa.Client, start, end time.Time, filter *audit.AuditLogFilter, chunkIndex int, totalChunks int, outputFile string, writeToStdout bool) ([]audit.AuditLog, error) {
	if !writeToStdout && totalChunks > 1 {
		fmt.Printf("Exporting chunk %d/%d: %s to %s...\n", chunkIndex+1, totalChunks, start.Format(time.RFC3339), end.Format(time.RFC3339))
	}

	// Fetch all audit logs with pagination
	allLogs := []audit.AuditLog{}
	page := 1
	limit := pageSize // Fetch logs per page

	for {
		request := &audit.ListAuditLogsRequest{
			Page:   page,
			Limit:  limit,
			Filter: filter,
		}

		result, err := client.Audit().ListAuditLogs(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch audit logs: %w", err)
		}
		fmt.Println("Fetched page", page, "items", len(result.Items), "total_pages", result.TotalPages, "total_items", result.TotalItems, "total_count", result.TotalCount)

		// Filter logs by time range
		for _, log := range result.Items {
			logTime := log.CreatedAt
			if !logTime.Before(start) && !logTime.After(end) {
				allLogs = append(allLogs, log)
			}
		}

		// Check if we've fetched all pages
		if page >= result.TotalPages || len(result.Items) == 0 {
			break
		}
		page++

		// Show progress (only if not writing to stdout and single file)
		if !writeToStdout && totalChunks == 1 && page%10 == 0 {
			fmt.Printf("Fetched %d pages, %d logs so far...\n", page-1, len(allLogs))
		}
	}

	if !writeToStdout && totalChunks > 1 {
		fmt.Printf("Found %d audit logs in chunk %d/%d\n", len(allLogs), chunkIndex+1, totalChunks)
	}

	return allLogs, nil
}

// writeExport writes the export data to a file
func writeExport(exportData map[string]interface{}, outputPath string, writeToStdout bool) error {
	var file *os.File
	var err error

	if writeToStdout {
		file = os.Stdout
	} else {
		// Ensure output directory exists
		outputDir := filepath.Dir(outputPath)
		if outputDir != "" && outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		file, err = os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"download"},
	Short:   "Export organisation audit logs to a JSON file",
	Long: `Export organisation audit logs to a JSON file for compliance purposes.

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
  tcloud audit export --since 1d --output -`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Parse API timeout (per chunk)
		timeoutDuration := 5 * time.Minute // default
		if apiTimeout != "" {
			parsedTimeout, err := time.ParseDuration(apiTimeout)
			if err != nil {
				return fmt.Errorf("invalid timeout duration: %w", err)
			}
			timeoutDuration = parsedTimeout
		}

		// Determine time range
		var start, end time.Time
		var rangeType string
		now := time.Now()

		// Determine split period
		if daily {
			splitPeriod = "daily"
		} else if weekly {
			splitPeriod = "weekly"
		} else if monthly {
			splitPeriod = "monthly"
		}

		// Determine time range from flags
		rangeFlagsSet := 0
		if sinceDuration != "" {
			rangeFlagsSet++
		}
		if startDate != "" || endDate != "" {
			rangeFlagsSet++
		}

		if rangeFlagsSet == 0 {
			return fmt.Errorf("must specify a time range: --since <duration> or --from/--to <dates>")
		}
		if rangeFlagsSet > 1 {
			return fmt.Errorf("cannot use --since with --from/--to, use one method")
		}

		if sinceDuration != "" {
			duration, err := parseDuration(sinceDuration)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}
			end = now
			start = now.Add(-duration)
			rangeType = "duration"
		} else {
			// Custom range with --from/--to
			rangeType = "custom"
			if startDate == "" {
				return fmt.Errorf("--from is required for custom time range")
			}
			if endDate == "" {
				return fmt.Errorf("--to is required for custom time range")
			}

			var parseErr error
			start, parseErr = time.Parse("2006-01-02", startDate)
			if parseErr != nil {
				return fmt.Errorf("invalid start date format (expected YYYY-MM-DD): %w", parseErr)
			}
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

			end, parseErr = time.Parse("2006-01-02", endDate)
			if parseErr != nil {
				return fmt.Errorf("invalid end date format (expected YYYY-MM-DD): %w", parseErr)
			}
			// Set end date to end of day
			end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		}

		if end.Before(start) {
			return fmt.Errorf("end date must be after start date")
		}

		writeToStdout := outputFile == "-"
		if !writeToStdout {
			fmt.Printf("Exporting audit logs from %s to %s...\n", start.Format(time.RFC3339), end.Format(time.RFC3339))
		}

		// Build filter if any filter flags are set
		var filter *audit.AuditLogFilter
		hasFilters := searchText != "" ||
			serviceAccount != "" ||
			userIdentity != "" ||
			impersonatorIdentity != "" ||
			len(actions) > 0 ||
			len(resourceTypes) > 0 ||
			resourceIdentity != "" ||
			organisationIdentity != "" ||
			includeSystemServices ||
			responseStatus != 0

		if hasFilters {
			filter = &audit.AuditLogFilter{
				SearchText:            searchText,
				ServiceAccount:        serviceAccount,
				UserIdentity:          userIdentity,
				ImpersonatorIdentity:  impersonatorIdentity,
				Actions:               actions,
				ResourceTypes:         resourceTypes,
				ResourceIdentity:      resourceIdentity,
				OrganizationIdentity:  organisationIdentity,
				IncludeSystemServices: includeSystemServices,
				ResponseStatus:        responseStatus,
			}
		}

		// Split time range if split period is specified
		chunks := splitTimeRange(start, end, splitPeriod)
		if len(chunks) == 0 {
			return fmt.Errorf("no time chunks to export")
		}

		// Export each chunk
		totalExported := 0
		for i, chunk := range chunks {
			chunkStart, chunkEnd := chunk[0], chunk[1]

			// Create context with timeout for this chunk
			chunkCtx, cancel := context.WithTimeout(cmd.Context(), timeoutDuration)

			// Fetch logs for this chunk
			chunkLogs, err := exportTimeRange(chunkCtx, client, chunkStart, chunkEnd, filter, i, len(chunks), outputFile, writeToStdout)

			// Always cancel the context when done with this chunk
			cancel()

			if err != nil {
				// Check if error is due to timeout
				if err == context.DeadlineExceeded || chunkCtx.Err() == context.DeadlineExceeded {
					return fmt.Errorf("timeout exceeded while exporting chunk %d/%d (%s to %s): %w", i+1, len(chunks), chunkStart.Format(time.RFC3339), chunkEnd.Format(time.RFC3339), err)
				}
				return err
			}

			// Determine output file for this chunk
			var chunkOutputFile string
			if writeToStdout && len(chunks) > 1 {
				return fmt.Errorf("cannot use stdout with split exports (--daily/--weekly/--monthly), specify an output file")
			}

			if writeToStdout {
				chunkOutputFile = "-"
			} else if len(chunks) > 1 {
				// Generate filename for this chunk
				if outputFile == "" {
					// Generate default filename with chunk date
					chunkOutputFile = fmt.Sprintf("audit-logs-%s-%s.json",
						chunkStart.Format("20060102"),
						chunkEnd.Format("20060102"))
				} else {
					// Use output file as template
					ext := filepath.Ext(outputFile)
					base := strings.TrimSuffix(outputFile, ext)
					chunkOutputFile = fmt.Sprintf("%s-%s-%s%s",
						base,
						chunkStart.Format("20060102"),
						chunkEnd.Format("20060102"),
						ext)
				}
			} else {
				// Single file export
				if outputFile == "" {
					timestamp := time.Now().Format("20060102-150405")
					chunkOutputFile = fmt.Sprintf("audit-logs-%s-%s.json", rangeType, timestamp)
				} else {
					chunkOutputFile = outputFile
				}
			}

			// Prepare export data with metadata
			exportData := map[string]interface{}{
				"exportedAt": time.Now().Format(time.RFC3339),
				"timeRange": map[string]interface{}{
					"type":  rangeType,
					"start": chunkStart.Format(time.RFC3339),
					"end":   chunkEnd.Format(time.RFC3339),
				},
				"totalLogs": len(chunkLogs),
				"logs":      chunkLogs,
			}

			if len(chunks) > 1 {
				exportData["chunk"] = map[string]interface{}{
					"index": i + 1,
					"total": len(chunks),
				}
			}

			// Add custom range dates if applicable
			if rangeType == "custom" {
				exportData["timeRange"].(map[string]interface{})["startDate"] = startDate
				exportData["timeRange"].(map[string]interface{})["endDate"] = endDate
			}

			// Add duration if using --since
			if rangeType == "duration" {
				exportData["timeRange"].(map[string]interface{})["since"] = sinceDuration
			}

			// Add filter metadata if filters were applied
			if hasFilters {
				filterMetadata := map[string]interface{}{}
				if searchText != "" {
					filterMetadata["searchText"] = searchText
				}
				if serviceAccount != "" {
					filterMetadata["serviceAccount"] = serviceAccount
				}
				if userIdentity != "" {
					filterMetadata["userIdentity"] = userIdentity
				}
				if impersonatorIdentity != "" {
					filterMetadata["impersonatorIdentity"] = impersonatorIdentity
				}
				if len(actions) > 0 {
					filterMetadata["actions"] = actions
				}
				if len(resourceTypes) > 0 {
					filterMetadata["resourceTypes"] = resourceTypes
				}
				if resourceIdentity != "" {
					filterMetadata["resourceIdentity"] = resourceIdentity
				}
				if organisationIdentity != "" {
					filterMetadata["organizationIdentity"] = organisationIdentity
				}
				if includeSystemServices {
					filterMetadata["includeSystemServices"] = includeSystemServices
				}
				if responseStatus != 0 {
					filterMetadata["responseStatus"] = responseStatus
				}
				exportData["filters"] = filterMetadata
			}

			// Write export
			if err := writeExport(exportData, chunkOutputFile, writeToStdout); err != nil {
				return err
			}

			if !writeToStdout {
				if len(chunks) > 1 {
					fmt.Fprintf(os.Stderr, "Exported %d audit logs to %s\n", len(chunkLogs), chunkOutputFile)
				} else {
					fmt.Fprintf(os.Stderr, "Successfully exported %d audit logs to %s\n", len(chunkLogs), chunkOutputFile)
				}
			}

			totalExported += len(chunkLogs)
		}

		if !writeToStdout && len(chunks) > 1 {
			fmt.Fprintf(os.Stderr, "Total: Exported %d audit logs across %d files\n", totalExported, len(chunks))
		}

		return nil
	},
}

func init() {
	AuditCmd.AddCommand(exportCmd)

	// Time range flags
	exportCmd.Flags().StringVar(&sinceDuration, "since", "", "Export logs from the past duration (e.g., 364d, 1w, 1mo, 1y, 24h)")
	exportCmd.Flags().StringVar(&startDate, "from", "", "Start date for custom range (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&endDate, "to", "", "End date for custom range (YYYY-MM-DD)")

	// Split flags
	exportCmd.Flags().BoolVar(&daily, "daily", false, "Split export into separate files per day")
	exportCmd.Flags().BoolVar(&weekly, "weekly", false, "Split export into separate files per week")
	exportCmd.Flags().BoolVar(&monthly, "monthly", false, "Split export into separate files per month")

	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (use '-' for stdout, default: audit-logs-{range}-{timestamp}.json)")

	// API timeout flag
	exportCmd.Flags().StringVar(&apiTimeout, "chunk-download-timeout", "", "API call timeout per chunk (e.g., 5m, 10m, 1h, default: 5m)")
	exportCmd.Flags().IntVar(&pageSize, "page-size", 100, "Page size for audit log list (default: 100)")

	// Filter flags
	exportCmd.Flags().StringVar(&searchText, "search-text", "", "Search text filter")
	exportCmd.Flags().StringVar(&serviceAccount, "service-account", "", "Filter by service account identity")
	exportCmd.Flags().StringVar(&userIdentity, "user-identity", "", "Filter by user identity")
	exportCmd.Flags().StringVar(&impersonatorIdentity, "impersonator-identity", "", "Filter by impersonator identity")
	exportCmd.Flags().StringSliceVar(&actions, "action", []string{}, "Filter by action(s) (can be specified multiple times)")
	exportCmd.Flags().StringSliceVar(&resourceTypes, "resource-type", []string{}, "Filter by resource type(s) (can be specified multiple times)")
	exportCmd.Flags().StringVar(&resourceIdentity, "resource-identity", "", "Filter by resource identity")
	exportCmd.Flags().StringVar(&organisationIdentity, "organisation-identity", "", "Filter by organisation identity")
	exportCmd.Flags().BoolVar(&includeSystemServices, "include-system-services", false, "Include system service logs")
	exportCmd.Flags().IntVar(&responseStatus, "response-status", 0, "Filter by HTTP response status code")
}
