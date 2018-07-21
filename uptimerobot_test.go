package uptimerobot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
    "github.com/stretchr/testify/assert"
)

func _makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Uptimerobot")
	}))
}

func TestMakeMonitor(t *testing.T) {
	ts := _makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		httpClient: ts.Client(),
		api_key:    "abcdefg",
	}
    m, err := c.getMonitors()
    assert.Equal(t, nil, err)
    assert.Equal(t, 0, len(m))
}
