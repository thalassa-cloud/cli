package workloadidentityfederation

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

func stdoutIsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func termGreen(s string) string {
	if !stdoutIsTTY() {
		return s
	}
	return "\x1b[32m" + s + "\x1b[0m"
}

func termBold(s string) string {
	if !stdoutIsTTY() {
		return s
	}
	return "\x1b[1m" + s + "\x1b[0m"
}

func termCyan(s string) string {
	if !stdoutIsTTY() {
		return s
	}
	return "\x1b[36m" + s + "\x1b[0m"
}

func termDim(s string) string {
	if !stdoutIsTTY() {
		return s
	}
	return "\x1b[2m" + s + "\x1b[0m"
}

func printBootstrapOutcome(vcs string, res *BootstrapResult, dry bool) {
	check := termGreen("✔")
	would := "○"
	if stdoutIsTTY() {
		would = "\x1b[33m○\x1b[0m" // amber for planned
	}

	fmt.Printf("%s Bootstrap workload identity (%s)\n\n", termCyan("►"), termBold(vcs))

	fmt.Printf("%s Organisation role - %s %s\n", check, res.RoleSlug, termDim("("+res.RoleIdentity+")"))

	// Federated identity provider (OIDC issuer registration)
	switch {
	case dry && res.WouldCreateProvider:
		fmt.Printf("%s Federated identity provider - would create %s\n", would, termDim("("+res.Issuer+")"))
	case dry:
		fmt.Printf("%s Federated identity provider - already present %s\n", check, termDim(res.ProviderIdentity))
	case res.CreatedProvider:
		fmt.Printf("%s Federated identity provider - created %s\n", check, res.ProviderIdentity)
	default:
		fmt.Printf("%s Federated identity provider - already present %s\n", check, res.ProviderIdentity)
	}

	// Thalassa service account
	switch {
	case dry && res.WouldCreateServiceAccount:
		fmt.Printf("%s Service account - would create\n", would)
	case dry:
		fmt.Printf("%s Service account - already present %s %s\n", check, res.ServiceAccountIdentity, termDim("("+res.ServiceAccountSlug+")"))
	case res.CreatedServiceAccount:
		fmt.Printf("%s Service account - created %s %s\n", check, res.ServiceAccountIdentity, termDim("("+res.ServiceAccountSlug+")"))
	default:
		fmt.Printf("%s Service account - already present %s %s\n", check, res.ServiceAccountIdentity, termDim("("+res.ServiceAccountSlug+")"))
	}

	// Federated identity (JWT subject → service account)
	switch {
	case dry && res.WouldCreateFederatedIdentity:
		fmt.Printf("%s Federated identity - would create %s\n", would, termDim("("+res.ProviderSubject+")"))
	case dry && res.WouldUpdateFederatedIdentity:
		fmt.Printf("%s Federated identity - would update configuration %s\n", would, res.FederatedIdentityIdentity)
	case dry:
		fmt.Printf("%s Federated identity - already present %s\n", check, res.FederatedIdentityIdentity)
	case res.CreatedFederatedIdentity:
		fmt.Printf("%s Federated identity - created %s\n", check, res.FederatedIdentityIdentity)
	case res.UpdatedFederatedIdentity:
		fmt.Printf("%s Federated identity - updated configuration %s\n", check, res.FederatedIdentityIdentity)
	default:
		fmt.Printf("%s Federated identity - already present %s\n", check, res.FederatedIdentityIdentity)
	}

	// Organisation role binding
	switch {
	case dry && res.WouldCreateRoleBinding:
		fmt.Printf("%s Organisation role binding - would create\n", would)
	case dry:
		fmt.Printf("%s Organisation role binding - already present\n", check)
	case res.CreatedRoleBinding:
		fmt.Printf("%s Organisation role binding - created\n", check)
	default:
		fmt.Printf("%s Organisation role binding - already present\n", check)
	}

	fmt.Println()
	fmt.Printf("  %s %s\n", termDim("issuer:"), res.Issuer)
	fmt.Printf("  %s %s\n", termDim("JWT sub:"), res.ProviderSubject)
}
