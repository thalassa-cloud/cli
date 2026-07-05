package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/config/securefile"
)

var migrateCredentials bool

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix config file security issues",
	Long:  "Fix security issues in the CLI config file, such as overly permissive file permissions or migrating credentials to the keychain.",
	Example: `  tcloud context fix
  tcloud context fix --migrate-credentials`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if migrateCredentials {
			if err := contextstate.MigrateCredentialsToKeychain(); err != nil {
				return err
			}
			fmt.Println("Migrated credentials from the config file to the keychain.")
		}

		filename := contextstate.ConfigFilename()
		mode, permissive, err := securefile.CheckPermissions(filename)
		if err != nil {
			return err
		}
		if !permissive {
			fmt.Printf("Config file %s already has secure permissions (%#o)\n", filename, mode)
			return nil
		}

		if err := contextstate.FixConfigPermissions(); err != nil {
			return err
		}
		fmt.Printf("Fixed permissions on config file %s (%#o -> %#o)\n", filename, mode, securefile.PrivateFileMode)
		return nil
	},
}

func init() {
	ContextCmd.AddCommand(fixCmd)
	fixCmd.Flags().BoolVar(&migrateCredentials, "migrate-credentials", false, "move plaintext credentials from the config file into the keychain")
}
