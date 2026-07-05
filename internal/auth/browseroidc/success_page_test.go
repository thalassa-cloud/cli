package browseroidc

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteLoginSuccessPage(t *testing.T) {
	rec := httptest.NewRecorder()
	writeLoginSuccessPage(rec)

	require.Equal(t, 200, rec.Code)
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Contains(t, body, "Signed in successfully")
	assert.Contains(t, body, "Thalassa Cloud")
	assert.Contains(t, body, "https://docs.thalassa.cloud/")
	assert.Contains(t, body, "https://thalassa.cloud/")
	assert.Contains(t, body, "<!DOCTYPE html>")
}
