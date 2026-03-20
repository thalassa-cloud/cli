package workloadidentityfederation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

func TestLabelsMatch(t *testing.T) {
	want := map[string]string{"a": "1", "b": "2"}
	assert.True(t, labelsMatch(map[string]string{"a": "1", "b": "2", "c": "3"}, want))
	assert.False(t, labelsMatch(map[string]string{"a": "1"}, want))
	assert.False(t, labelsMatch(map[string]string{"a": "1", "b": "x"}, want))
	assert.False(t, labelsMatch(nil, want))
	assert.False(t, labelsMatch(map[string]string{}, map[string]string{}))
}

func TestIssuerURLHostname(t *testing.T) {
	assert.Equal(t, "gitlab.com", issuerURLHostname("https://gitlab.com"))
	assert.Equal(t, "gitlab.com", issuerURLHostname("https://gitlab.com/"))
	assert.Equal(t, "gitlab.example.org", issuerURLHostname("https://gitlab.example.org:8443"))
	assert.Equal(t, "k8s.example.com", issuerURLHostname("k8s.example.com"))
	assert.Equal(t, "", issuerURLHostname(""))
}

func TestWifResourceKeyIncludesIssuer(t *testing.T) {
	subject := "project_path:g/p:ref_type:branch:ref:main"
	a := wifResourceKey(ValueVCSGitLab, "g/p", subject, "https://gitlab.com")
	b := wifResourceKey(ValueVCSGitLab, "g/p", subject, "https://gitlab.example.com")
	assert.NotEqual(t, a, b, "different issuers must not share the same wif-key")
	c := wifResourceKey(ValueVCSGitLab, "g/p", subject, "https://gitlab.com/")
	assert.Equal(t, a, c, "issuer normalization should match")
}

func TestFederatedIdentityNeedsBootstrapReconcile(t *testing.T) {
	key := "abc123"
	fi := &clientiam.FederatedIdentity{
		Name:              "n",
		Description:       "d",
		AllowedScopes:     []clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIRead},
		TrustedAudiences:  []string{"https://api.example"},
		AudienceMatchMode: clientiam.AudienceMatchModeAny,
		Labels:            bootstrapLabels(ValueVCSGitLab, key),
		Annotations:       map[string]string{"thalassa.cloud/wif.provider-subject": "sub"},
	}
	wantScopes := []clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIRead}
	wantAud := []string{"https://api.example"}
	assert.False(t, federatedIdentityNeedsBootstrapReconcile(fi, "n", "d", "sub", ValueVCSGitLab, key, wantScopes, wantAud))
	assert.True(t, federatedIdentityNeedsBootstrapReconcile(fi, "n", "d", "sub", ValueVCSGitLab, key,
		[]clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIRead, clientiam.AccessCredentialsScopeAPIWrite}, wantAud))
}

func TestScopesEqual(t *testing.T) {
	a := []clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIRead, clientiam.AccessCredentialsScopeAPIWrite}
	b := []clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIWrite, clientiam.AccessCredentialsScopeAPIRead}
	c := []clientiam.AccessCredentialsScope{clientiam.AccessCredentialsScopeAPIRead}
	assert.True(t, scopesEqual(a, b))
	assert.True(t, scopesEqual(nil, nil))
	assert.True(t, scopesEqual(nil, []clientiam.AccessCredentialsScope{}))
	assert.False(t, scopesEqual(a, c))
}
