package securitygroups

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var RulesCmd = &cobra.Command{
	Use:     "rules",
	Aliases: []string{"rule"},
	Short:   "Manage security group rules",
}

var (
	setIngressFile string
	setEgressFile  string
)

func loadRules(path string) ([]iaas.SecurityGroupRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file: %w", err)
	}
	var rules []iaas.SecurityGroupRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parse rules JSON: %w", err)
	}
	return rules, nil
}

var setIngressCmd = &cobra.Command{
	Use:               "set-ingress SECURITY_GROUP",
	Short:             "Replace all ingress rules from a JSON file",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeSecurityGroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := loadRules(setIngressFile)
		if err != nil {
			return err
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		result, err := client.IaaS().BatchUpdateSecurityGroupIngressRules(cmd.Context(), args[0], iaas.BatchUpdateSecurityGroupRulesRequest{Rules: rules})
		if err != nil {
			return fmt.Errorf("failed to set ingress rules: %w", err)
		}
		fmt.Printf("Set %d ingress rules on %s\n", len(result), args[0])
		return nil
	},
}

var setEgressCmd = &cobra.Command{
	Use:               "set-egress SECURITY_GROUP",
	Short:             "Replace all egress rules from a JSON file",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeSecurityGroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := loadRules(setEgressFile)
		if err != nil {
			return err
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		result, err := client.IaaS().BatchUpdateSecurityGroupEgressRules(cmd.Context(), args[0], iaas.BatchUpdateSecurityGroupRulesRequest{Rules: rules})
		if err != nil {
			return fmt.Errorf("failed to set egress rules: %w", err)
		}
		fmt.Printf("Set %d egress rules on %s\n", len(result), args[0])
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(RulesCmd)

	RulesCmd.AddCommand(setIngressCmd)
	setIngressCmd.Flags().StringVar(&setIngressFile, "file", "", "JSON array of SecurityGroupRule objects")
	_ = setIngressCmd.MarkFlagRequired("file")

	RulesCmd.AddCommand(setEgressCmd)
	setEgressCmd.Flags().StringVar(&setEgressFile, "file", "", "JSON array of SecurityGroupRule objects")
	_ = setEgressCmd.MarkFlagRequired("file")
}
