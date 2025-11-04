package oidc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

var (
	subjectTokenFlag        string
	organisationIDFlag      string
	serviceAccountIDFlag    string
	accessTokenLifetimeFlag string
)

// tokenExchangeCmd represents the token-exchange command
var tokenExchangeCmd = &cobra.Command{
	Use:   "token-exchange",
	Short: "Exchange an OIDC token for an access token",
	Long: "Helper for exchanging an OIDC subject token for a Thalassa Cloud access token using the token exchange grant type. " +
		"Intended for use in CI/CD pipelines such as GitLab CI, GitHub Actions, Kubernetes, and similar automation environments " +
		"where OIDC identity tokens are provided and need to be exchanged for a Thalassa access token.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get subject token from flag or environment variable
		subjectToken := subjectTokenFlag
		if subjectToken == "" {
			subjectToken = os.Getenv("THALASSA_SUBJECT_ID_TOKEN")
		}
		if subjectToken == "" {
			return errors.New("subject token is required (use --subject-token flag or set THALASSA_SUBJECT_ID_TOKEN environment variable)")
		}

		// Get organisation ID from flag or context
		organisationID := organisationIDFlag
		if organisationID == "" {
			organisationID = contextstate.Organisation()
		}
		if organisationID == "" {
			return errors.New("organisation ID is required (use --organisation-id flag or set organisation in context)")
		}

		// Get service account ID from flag or environment variable
		serviceAccountID := serviceAccountIDFlag
		if serviceAccountID == "" {
			serviceAccountID = os.Getenv("THALASSA_SERVICE_ACCOUNT_ID")
		}
		if serviceAccountID == "" {
			return errors.New("service account ID is required (use --service-account-id flag or set THALASSA_SERVICE_ACCOUNT_ID environment variable)")
		}

		// Get API endpoint from context or flag
		apiEndpoint := contextstate.Server()
		if apiEndpoint == "" {
			return errors.New("API endpoint is required (use --api flag or set in context)")
		}

		// Parse and validate access token lifetime
		accessTokenLifetime := accessTokenLifetimeFlag
		if accessTokenLifetime == "" {
			accessTokenLifetime = "1h"
		}

		duration, err := time.ParseDuration(accessTokenLifetime)
		if err != nil {
			return fmt.Errorf("invalid access token lifetime format: %w", err)
		}

		minDuration := time.Minute
		maxDuration := 24 * time.Hour
		if duration < minDuration {
			return fmt.Errorf("access token lifetime must be at least %s (got %s)", minDuration, duration)
		}
		if duration > maxDuration {
			return fmt.Errorf("access token lifetime must be at most %s (got %s)", maxDuration, duration)
		}

		// Build the token exchange URL
		tokenURL := fmt.Sprintf("%s/oidc/token", apiEndpoint)

		// Prepare form data
		formData := url.Values{}
		formData.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
		formData.Set("subject_token", subjectToken)
		formData.Set("subject_token_type", "urn:ietf:params:oauth:token-type:jwt")
		formData.Set("organisation_id", organisationID)
		formData.Set("service_account_id", serviceAccountID)
		formData.Set("access_token_lifetime", fmt.Sprintf("%d", int64(duration.Seconds())))

		// Create HTTP request
		req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		// Check for HTTP errors
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
		}

		// Parse JSON response
		var tokenResponse struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if tokenResponse.AccessToken == "" {
			return errors.New("access token not found in response")
		}

		// Output the access token
		fmt.Println(tokenResponse.AccessToken)
		return nil
	},
}

func init() {
	tokenExchangeCmd.Flags().StringVar(&subjectTokenFlag, "subject-token", "", "Subject token (JWT) to exchange (can also be set via THALASSA_ID_TOKEN env var)")
	tokenExchangeCmd.Flags().StringVar(&organisationIDFlag, "organisation-id", "", "Organisation ID (can also be set via context)")
	tokenExchangeCmd.Flags().StringVar(&serviceAccountIDFlag, "service-account-id", "", "Service account ID (can also be set via THALASSA_SERVICE_ACCOUNT_ID env var)")
	tokenExchangeCmd.Flags().StringVar(&accessTokenLifetimeFlag, "access-token-lifetime", "1h", "Access token lifetime (min: 1m, max: 24h, default: 1h)")
}
