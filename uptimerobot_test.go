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
        var mons = []Monitor {
            Monitor {
                Id: "0",
                Friendly_name: "foo",
                Url: "http://nothing.test",
                Monitor_type: 1,
            },
        }
        p := Pagination {
            Offset: 0,
            Limit: 1,
            Total: 1,
        }
        mr := MonitorResp {
            Stat: "ok",
            Pagination: p,
            Monitors: mons,
        }
        e.Encode(mr)
	}))
}

func TestGetMonitors(t *testing.T) {
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
    assert.Equal(t, 1, len(m))
    assert.Equal(t, "0", (m[0].Id))
}
