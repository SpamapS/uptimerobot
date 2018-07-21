package uptimerobot

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
	api_key    string
}

type Log struct {
	log_type int
	datetime int
	duration int
}

type Monitor struct {
	id              string
	friendly_name   string
	url             string
	monitor_type    int
	sub_type        *string
	keyword_type    *string
	keyword_value   *string
	http_username   *string
	http_password   *string
	port            *string
	interval        *int
	status          *int
	create_datetime *int
	monitor_group   *int
	is_group_main   *int
	logs            []Log
}

type Pagination struct {
	offset int
	limit  int
	total  int
}

type MonitorResp struct {
	stat       string
	pagination Pagination
	monitors   []Monitor
}

type CreatedMonitor struct {
	id     string
	status int
}

type CreateMonitorResp struct {
	stat    string
	monitor CreatedMonitor
}

func (c *Client) _makeReq(path string, data *url.Values) (*http.Request, error) {
	data.Set("api_key", c.api_key)
	data.Set("format", "json")
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) getMonitors() ([]Monitor, error) {
	data := url.Values{}

	req, err := c._makeReq("/getMonitors", &data)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var monitors_resp MonitorResp
	err = json.NewDecoder(resp.Body).Decode(&monitors_resp)
	if err != nil {
		return nil, err
	}
	return monitors_resp.monitors, err
}

func (c *Client) createMonitor(friendly_name string, monitor_url string, monitor_type string) (string, error) {
	data := url.Values{}
	req, err := c._makeReq("/newMonitor", &data)
	if err != nil {
		return "", err
	}
	data.Set("friendly_name", friendly_name)
	data.Set("url", monitor_url)
	data.Set("type", monitor_type)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var monitor_create_resp CreateMonitorResp
	err = json.NewDecoder(resp.Body).Decode(&monitor_create_resp)
	if err != nil {
		return "", err
	}
	return monitor_create_resp.monitor.id, nil
}
