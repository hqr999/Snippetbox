package main

import (
	"net/http"
	"testing"

	"github.com/hqr999/Snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	t_server := newTestServer(t, app.routes())
	defer t_server.Close()

	res := t_server.get(t, "/ping")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}
