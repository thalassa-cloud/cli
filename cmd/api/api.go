package api

import (
	"github.com/spf13/cobra"
)

// ApiCmd represents the api command
var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Direct API access",
	Long:  "Make raw HTTP requests to the Thalassa Cloud API.",
}
