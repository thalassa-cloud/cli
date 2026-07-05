package browseroidc

import (
	_ "embed"
	"net/http"
)

//go:embed success.html
var loginSuccessHTML []byte

func writeLoginSuccessPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(loginSuccessHTML)
}
