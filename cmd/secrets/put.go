package secrets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientsecrets "github.com/thalassa-cloud/client-go/secrets"
)

var (
	putRegion      string
	putPath        string
	putString      string
	putFromFile    string
	putKV          []string
	putGenerateLen int
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Put a new secret version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		secretString := putString
		if putFromFile != "" {
			data, err := os.ReadFile(putFromFile)
			if err != nil {
				return fmt.Errorf("read --from-file: %w", err)
			}
			secretString = string(data)
		}

		req := clientsecrets.PutSecretValueRequest{
			Path:            putPath,
			SecretString:    secretString,
			SecretKeyValues: shared.KeyValuePairsToMap(putKV),
		}
		if putGenerateLen > 0 {
			if err := validateGenerateBytes(putGenerateLen); err != nil {
				return err
			}
			req.GenerateSecret = &clientsecrets.GenerateSecret{ByteLength: putGenerateLen}
		}
		if req.SecretString == "" && len(req.SecretKeyValues) == 0 && req.GenerateSecret == nil {
			return fmt.Errorf("provide --string, --from-file, --kv, or --generate-bytes[=N] (bare --generate-bytes defaults to %d)", generateBytesDefault)
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		result, err := client.Secrets().PutSecretValue(cmd.Context(), putRegion, putPath, req)
		if err != nil {
			return fmt.Errorf("failed to put secret value: %w", err)
		}
		fmt.Printf("Secret %s updated to version %d\n", result.Path, result.Version)
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(putCmd)
	putCmd.Flags().StringVar(&putRegion, "region", "", "Region")
	putCmd.Flags().StringVar(&putPath, "path", "", "Secret path")
	putCmd.Flags().StringVar(&putString, "string", "", "Secret string value")
	putCmd.Flags().StringVar(&putFromFile, "from-file", "", "Read secret string from a file")
	putCmd.Flags().StringSliceVar(&putKV, "kv", nil, "Secret key/value pairs as key=value (repeatable)")
	putCmd.Flags().IntVar(&putGenerateLen, "generate-bytes", 0, generateBytesFlagUsage())
	putCmd.Flags().Lookup("generate-bytes").NoOptDefVal = fmt.Sprintf("%d", generateBytesDefault)
	_ = putCmd.MarkFlagRequired("region")
	_ = putCmd.MarkFlagRequired("path")
	_ = putCmd.MarkFlagFilename("from-file")
	_ = putCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = putCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
}
