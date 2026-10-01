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

func bootstrapOutcomeMarkers() (check, would string) {
	check = termGreen("✔")
	would = "○"
	if stdoutIsTTY() {
		would = "\x1b[33m○\x1b[0m" // amber for planned
	}
	return check, would
}

func printBootstrapOutcome(vcs string, res *BootstrapResult, dry bool) {
	check, would := bootstrapOutcomeMarkers()

	fmt.Printf("%s Bootstrap workload identity (%s)\n\n", termCyan("►"), termBold(vcs))

	for _, role := range res.Roles {
		fmt.Printf("%s Organisation role - %s %s\n", check, role.Slug, termDim("("+role.Identity+")"))
	}
	for _, policy := range res.Policies {
		fmt.Printf("%s IAM policy - %s %s\n", check, policy.Slug, termDim("("+policy.Identity+")"))
	}

	printBootstrapProviderOutcome(check, would, res, dry)
	printBootstrapServiceAccountOutcome(check, would, res, dry)
	printBootstrapFederatedIdentityOutcome(check, would, res, dry)
	printBootstrapRoleBindingOutcomes(check, would, res.Roles, dry)
	printBootstrapPolicyBindingOutcomes(check, would, res.Policies, dry)

	fmt.Println()
	fmt.Printf("  %s %s\n", termDim("issuer:"), res.Issuer)
	fmt.Printf("  %s %s\n", termDim("JWT sub:"), res.ProviderSubject)
}

func printBootstrapProviderOutcome(check, would string, res *BootstrapResult, dry bool) {
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
}

func printBootstrapServiceAccountOutcome(check, would string, res *BootstrapResult, dry bool) {
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
}

func printBootstrapFederatedIdentityOutcome(check, would string, res *BootstrapResult, dry bool) {
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
}

func printBootstrapBindingLine(check, would, kind, label string, dry, wouldCreate, created bool) {
	switch {
	case dry && wouldCreate:
		fmt.Printf("%s %s (%s) - would create\n", would, kind, label)
	case dry:
		fmt.Printf("%s %s (%s) - already present\n", check, kind, label)
	case created:
		fmt.Printf("%s %s (%s) - created\n", check, kind, label)
	default:
		fmt.Printf("%s %s (%s) - already present\n", check, kind, label)
	}
}

func printBootstrapRoleBindingOutcomes(check, would string, roles []BootstrapRoleResult, dry bool) {
	for _, role := range roles {
		label := role.Slug
		if label == "" {
			label = role.Identity
		}
		printBootstrapBindingLine(check, would, "Organisation role binding", label, dry, role.WouldCreateBinding, role.CreatedBinding)
	}
}

func printBootstrapPolicyBindingOutcomes(check, would string, policies []BootstrapPolicyResult, dry bool) {
	for _, policy := range policies {
		label := policy.Slug
		if label == "" {
			label = policy.Identity
		}
		printBootstrapBindingLine(check, would, "IAM policy binding", label, dry, policy.WouldCreateBinding, policy.CreatedBinding)
	}
}
