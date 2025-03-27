package me

import (
	"github.com/spf13/cobra"
)

var MeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get information about the current user",
}
