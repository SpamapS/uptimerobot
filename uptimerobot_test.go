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
        if r.Method == "POST" && r.URL.Path == "/getMonitors" {
            var mons = []Monitor {
                Monitor {
                    Id: 99,
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
        } else if r.Method == "POST" && r.URL.Path == "/newMonitor" {
            var created = CreateMonitorResp {
                Stat: "ok",
                Monitor: CreatedMonitor {
                    Id: 99,
                    Status: 1,
                },
            }
            e.Encode(created)
        }
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
    assert.Equal(t, 99, (m[0].Id))
}

func TestCreateMonitor(t *testing.T) {
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
    new_id, err := c.createMonitor("make_friendly", "http://make.test", MONITOR_TYPE_HTTP)
    assert.Equal(t, nil, err)
    assert.Equal(t, 99, new_id)
}
