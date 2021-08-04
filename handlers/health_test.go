package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHealthHandler(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/healthz", strings.NewReader("body"))
	if err != nil {
		t.Fatal(err)
	}

	handler := MakeHealthHandler()

	handler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
