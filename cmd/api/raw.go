package api

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	rawMethod      string
	rawData        string
	rawShowHeaders bool
)

// rawCmd represents the raw API request command
var rawCmd = &cobra.Command{
	Use:   "raw PATH",
	Short: "Make a raw HTTP request to the API",
	Long: `Make a raw HTTP request to the Thalassa Cloud API.

Similar to 'kubectl get --raw', this bypasses the CLI resource layer and sends
the request directly to the API server. Uses the same authentication and
context (organisation, endpoint) as other tcloud commands.

PATH must start with a slash (e.g. /v1/me/organisations).
Requires client-go with RawRequest support.`,
	Example: `  tcloud api raw /v1/me/organisations
  tcloud api raw -X GET /v1/iaas/regions
  tcloud api raw -X POST -d '{"name":"test"}' /v1/some/resource
  tcloud api raw --show-headers /v1/me`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("path must start with /")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		rq := client.GetClient()
		method := strings.ToUpper(rawMethod)
		if method == "" {
			method = "GET"
		}

		var body []byte
		if rawData != "" {
			body = []byte(rawData)
		}

		resp, err := rq.RawRequest(cmd.Context(), method, path, body)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.RawResponse.Body.Close()

		if rawShowHeaders {
			for k, v := range resp.RawResponse.Header {
				fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
			}
			fmt.Println()
		}

		out := resp.Body()
		fmt.Print(string(out))
		if len(out) > 0 && out[len(out)-1] != '\n' {
			fmt.Println()
		}

		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			return fmt.Errorf("API returned status %d", resp.StatusCode())
		}
		return nil
	},
}

func init() {
	ApiCmd.AddCommand(rawCmd)

	rawCmd.Flags().StringVarP(&rawMethod, "request", "X", "GET", "HTTP method (GET, POST, PUT, PATCH, DELETE)")
	rawCmd.Flags().StringVarP(&rawData, "data", "d", "", "Request body (for POST, PUT, PATCH)")
	rawCmd.Flags().BoolVar(&rawShowHeaders, "show-headers", false, "Print response headers")
}
