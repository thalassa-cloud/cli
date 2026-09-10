package secrets

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientsecrets "github.com/thalassa-cloud/client-go/secrets"
)

var (
	createRegion      string
	createPath        string
	createDescription string
	createKmsKey      string
	createString      string
	createFromFile    string
	createKV          []string
	createGenerateLen int
	createLabels      []string
	createAnnotations []string
	createPolicyFile  string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a secret (metadata response only; use get-value to read material)",
	Example: `  tcloud secrets create --region nl-01 --path /app/prod/db --kms-key kms-123 --generate-bytes
  tcloud secrets create --region nl-01 --path /app/prod/db --kms-key kms-123 --generate-bytes=64
  tcloud secrets create --region nl-01 --path /app/prod/token --kms-key kms-123 --from-file ./token.txt`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createPath == "" {
			return fmt.Errorf("--path is required")
		}
		if createKmsKey == "" {
			return fmt.Errorf("--kms-key is required")
		}
		if createString == "" && createFromFile == "" && len(createKV) == 0 && createGenerateLen <= 0 {
			return fmt.Errorf("provide --string, --from-file, --kv, or --generate-bytes[=N] (bare --generate-bytes defaults to %d)", generateBytesDefault)
		}
		if createGenerateLen > 0 {
			if err := validateGenerateBytes(createGenerateLen); err != nil {
				return err
			}
		}

		secretString := createString
		if createFromFile != "" {
			data, err := os.ReadFile(createFromFile)
			if err != nil {
				return fmt.Errorf("read --from-file: %w", err)
			}
			secretString = string(data)
		}

		req := clientsecrets.CreateSecretRequest{
			Path:            createPath,
			Description:     createDescription,
			Labels:          shared.KeyValuePairsToMap(createLabels),
			Annotations:     shared.KeyValuePairsToMap(createAnnotations),
			KmsKeyIdentity:  createKmsKey,
			SecretString:    secretString,
			SecretKeyValues: shared.KeyValuePairsToMap(createKV),
		}
		if createGenerateLen > 0 {
			req.GenerateSecret = &clientsecrets.GenerateSecret{ByteLength: createGenerateLen}
		}
		if createPolicyFile != "" {
			policy, err := readSecretPolicy(createPolicyFile)
			if err != nil {
				return err
			}
			req.AccessPolicy = policy
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		secret, err := client.Secrets().CreateSecret(cmd.Context(), createRegion, req)
		if err != nil {
			return fmt.Errorf("failed to create secret: %w", err)
		}

		kmsKey := "-"
		if secret.KmsKey != nil {
			kmsKey = secret.KmsKey.Identity
		}
		body := [][]string{{
			secret.Path,
			fmt.Sprintf("%d", secret.CurrentVersion),
			kmsKey,
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Path", "Version", "KMS Key"}, body)
		}
		return nil
	},
}

func readSecretPolicy(path string) (*clientsecrets.SecretPolicy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy file: %w", err)
	}
	var policy clientsecrets.SecretPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("parse policy JSON: %w", err)
	}
	return &policy, nil
}

func init() {
	SecretsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region")
	createCmd.Flags().StringVar(&createPath, "path", "", "Secret path")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringVar(&createKmsKey, "kms-key", "", "KMS key identity used to encrypt the secret")
	createCmd.Flags().StringVar(&createString, "string", "", "Secret string value")
	createCmd.Flags().StringVar(&createFromFile, "from-file", "", "Read secret string from a file")
	createCmd.Flags().StringSliceVar(&createKV, "kv", nil, "Secret key/value pairs as key=value (repeatable)")
	createCmd.Flags().IntVar(&createGenerateLen, "generate-bytes", 0, generateBytesFlagUsage())
	createCmd.Flags().Lookup("generate-bytes").NoOptDefVal = fmt.Sprintf("%d", generateBytesDefault)
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createPolicyFile, "policy-file", "", "JSON file with an access policy")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.MarkFlagRequired("path")
	_ = createCmd.MarkFlagRequired("kms-key")
	_ = createCmd.MarkFlagFilename("from-file")
	_ = createCmd.MarkFlagFilename("policy-file")
	_ = createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = createCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
	_ = createCmd.RegisterFlagCompletionFunc("kms-key", completion.CompleteKmsKeyIdentity)
}
