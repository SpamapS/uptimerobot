package uptimerobot

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := json.NewEncoder(w)
		if r.Method == "POST" && r.URL.Path == "/getMonitors" {
			r.ParseForm()
			mv := r.PostFormValue("monitors")
			var mons []Monitor
			mon_99 := Monitor{
				Id:            99,
				Friendly_name: "foo",
				Url:           "http://nothing.test",
				Monitor_type:  1,
			}
			mon_100 := Monitor{
				Id:            100,
				Friendly_name: "bar",
				Url:           "http://nobar.test",
				Monitor_type:  1,
			}
			if mv == "" || mv == "99-100" || mv == "100-99" {
				mons = []Monitor{
					mon_99,
					mon_100,
				}
			} else if mv == "99" {
				mons = []Monitor{
					mon_99,
				}
			} else if mv == "100" {
				mons = []Monitor{
					mon_100,
				}
			}
			p := Pagination{
				Offset: 0,
				Limit:  2,
				Total:  len(mons),
			}
			mr := MonitorResp{
				Stat:       "ok",
				Pagination: p,
				Monitors:   mons,
			}
			e.Encode(mr)
		} else if r.Method == "POST" && r.URL.Path == "/newMonitor" {
			var status = 1
			var created = ChangeMonitorResp{
				Stat: "ok",
				Monitor: ChangedMonitor{
					Id:     99,
					Status: &status,
				},
			}
			e.Encode(created)
		} else if r.Method == "POST" && r.URL.Path == "/editMonitor" {
			var edited = ChangeMonitorResp{
				Stat: "ok",
				Monitor: ChangedMonitor{
					Id: 99,
				},
			}
			e.Encode(edited)
		} else if r.Method == "POST" && r.URL.Path == "/deleteMonitor" {
			var deleted = ChangeMonitorResp{
				Stat: "ok",
				Monitor: ChangedMonitor{
					Id: 99,
				},
			}
			e.Encode(deleted)
		} else {
			http.NotFound(w, r)
		}
	}))
}

func TestGetMonitors(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		HttpClient: ts.Client(),
		Api_key:    "abcdefg",
	}
	var ids []int
	m, err := c.GetMonitors(ids)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 99, (m[0].Id))
}

func TestGetMonitorsIds(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		HttpClient: ts.Client(),
		Api_key:    "abcdefg",
	}
	var ids []int = []int{99}
	m, err := c.GetMonitors(ids)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 99, (m[0].Id))
}

func TestCreateMonitor(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		HttpClient: ts.Client(),
		Api_key:    "abcdefg",
	}
	m := Monitor{
		Friendly_name: "make_friendly",
		Url:           "http://make.test",
		Monitor_type:  MONITOR_TYPE_HTTP,
	}
	err = c.CreateMonitor(&m)
	assert.Equal(t, nil, err)
	assert.Equal(t, 99, m.Id)
}

func TestEditMonitor(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		HttpClient: ts.Client(),
		Api_key:    "abcdefg",
	}
	m := Monitor{
		Id:            99,
		Friendly_name: "make_friendly",
		Url:           "http://make.test",
		Monitor_type:  MONITOR_TYPE_HTTP,
	}
	err = c.EditMonitor(&m)
	assert.Equal(t, nil, err)
}

func TestDeleteMonitor(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Equal(t, nil, err)

	c := Client{
		BaseURL:    u,
		UserAgent:  "Bah",
		HttpClient: ts.Client(),
		Api_key:    "abcdefg",
	}
	err = c.DeleteMonitor(99)
	assert.Equal(t, nil, err)
}
