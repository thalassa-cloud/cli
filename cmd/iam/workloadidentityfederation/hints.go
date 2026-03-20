package workloadidentityfederation

import (
	"fmt"
	"strings"
)

func printBootstrapHints(opts BootstrapOptions, res *BootstrapResult, apiBase string, organisation string) {
	tokenURL := apiBase + "/oidc/token"
	primaryAud := opts.TrustedAudiences[0]
	if len(opts.TrustedAudiences) > 1 {
		fmt.Printf("\nNote: federated identity trusts multiple audiences; CI id_token aud must be one of: %s\n",
			strings.Join(opts.TrustedAudiences, ", "))
	}

	fmt.Println()
	fmt.Printf("%s Next steps CI configuration (https://docs.thalassa.cloud/docs/iam/oidc/) \n", termCyan("►"))
	switch opts.VCS {
	case ValueVCSGitHub:
		fmt.Printf(`
name: Deploy to Thalassa Cloud
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
    steps:
      - uses: actions/checkout@v4
      - name: Get OIDC token
        id: oidc
        uses: actions/github-script@v7
        with:
          script: |
            const token = await core.getIDToken(%q);
            core.setOutput('token', token);
      - name: Exchange OIDC token
        id: auth
        run: |
          RESPONSE=$(curl -sS -X POST %s \
            -H "Content-Type: application/x-www-form-urlencoded" \
            -d "grant_type=urn:ietf:params:oauth:grant-type:token-exchange" \
            -d "subject_token=${{ steps.oidc.outputs.token }}" \
            -d "subject_token_type=urn:ietf:params:oauth:token-type:id_token" \
            -d "organisation_id=${{ vars.THALASSA_ORGANISATION_ID }}" \
            -d "service_account_id=${{ vars.THALASSA_SERVICE_ACCOUNT_ID }}")
          echo "token=$(echo "$RESPONSE" | jq -r '.access_token')" >> "$GITHUB_OUTPUT"
      - name: Example use
        env:
          THALASSA_TOKEN: ${{ steps.auth.outputs.token }}
        run: echo "THALASSA_TOKEN is set"
`, primaryAud, tokenURL)

	case ValueVCSGitLab:
		fmt.Printf(`
stages:
  - deploy

deploy:
  stage: deploy
  image: alpine:latest
  id_tokens:
    THALASSA_ID_TOKEN:
      aud: %q
  before_script:
    - apk add --no-cache curl jq
  variables:
    THALASSA_SERVICE_ACCOUNT_ID: "%s"
    THALASSA_ORGANISATION_ID: "<your-organisation-id>"
  script:
    - |
      BEARER_TOKEN=$(curl -sS -X POST %s \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=urn:ietf:params:oauth:grant-type:token-exchange" \
        -d "subject_token=${THALASSA_ID_TOKEN}" \
        -d "subject_token_type=urn:ietf:params:oauth:token-type:id_token" \
        -d "organisation_id=${THALASSA_ORGANISATION_ID}" \
        -d "service_account_id=${THALASSA_SERVICE_ACCOUNT_ID}" \
        | jq -r '.access_token')
    - export THALASSA_TOKEN="${BEARER_TOKEN}"
`, primaryAud, res.ServiceAccountIdentity, tokenURL)

	case ValueVCSKubernetes:
		fmt.Printf(`
Kubernetes: request a projected service account token whose audience matches the federated identity, then exchange it.

apiVersion: v1
kind: Pod
metadata:
  name: thalassa-example
spec:
  serviceAccountName: <your-namespace-sa>
  volumes:
    - name: thalassa-oidc
      projected:
        sources:
          - serviceAccountToken:
              path: thalassa-oidc
              expirationSeconds: 3600
              audience: %q
  containers:
    - name: app
      image: alpine:latest
      volumeMounts:
        - name: thalassa-oidc
          mountPath: /var/run/secrets/thalassa
          readOnly: true
      env:
        - name: THALASSA_SERVICE_ACCOUNT_ID
          value: "%s"
        - name: THALASSA_ORGANISATION_ID
          value: "<your-organisation-id>"
      command: ["/bin/sh", "-c"]
      args:
        - |
          apk add --no-cache curl jq
          export THALASSA_SUBJECT_ID_TOKEN="$(cat /var/run/secrets/thalassa/thalassa-oidc)"
          BEARER_TOKEN=$(curl -sS -X POST %s \
            -H "Content-Type: application/x-www-form-urlencoded" \
            -d "grant_type=urn:ietf:params:oauth:grant-type:token-exchange" \
            -d "subject_token=${THALASSA_SUBJECT_ID_TOKEN}" -d "subject_token_type=urn:ietf:params:oauth:token-type:id_token" \
            -d "organisation_id=${THALASSA_ORGANISATION_ID}" -d "service_account_id=${THALASSA_SERVICE_ACCOUNT_ID}" \
            | jq -r '.access_token')
          export THALASSA_TOKEN="${BEARER_TOKEN}"
          # ... call Thalassa API or run tcloud with THALASSA_TOKEN
`, primaryAud, res.ServiceAccountIdentity, tokenURL)
	}
}
