package uptimerobot

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
    "github.com/stretchr/testify/assert"
    "encoding/json"
)

func _makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        e := json.NewEncoder(w)
        cmon := CreateMonitorResp {
            stat: "ok",
            monitor: CreatedMonitor {
                id: "0",
                status: 0,
            },
        }
        e.Encode(cmon)
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
    assert.Equal(t, "ok", m.stat)
}
