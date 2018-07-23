package uptimerobot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HttpClient *http.Client
	Api_key    string
}

type Log struct {
	Log_type int `json:"type"`
	Datetime int `json:"datetime"`
	Duration int `json:"duration"`
}

const MONITOR_TYPE_HTTP = 1
const MONITOR_TYPE_KEYWORD = 2
const MONITOR_TYPE_PING = 3
const MONITOR_TYPE_PORT = 4

const MONITOR_SUB_TYPE_HTTP = 1
const MONITOR_SUB_TYPE_HTTPS = 2
const MONITOR_SUB_TYPE_FTP = 3
const MONITOR_SUB_TYPE_SMTP = 4
const MONITOR_SUB_TYPE_POP3 = 5
const MONITOR_SUB_TYPE_IMAP = 6
const MONITOR_SUB_TYPE_CUSTOM = 99

const MONITOR_STATUS_PAUSED = 0
const MONITOR_STATUS_NOT_CHECKED = 1
const MONITOR_STATUS_UP = 2
const MONITOR_STATUS_SEEMS_DOWN = 3
const MONITOR_STATUS_DOWN = 4

const KEYWORD_TYPE_EXISTS = 1
const KEYWORD_TYPE_NOT_EXISTS = 2

type Monitor struct {
	Id              int     `json:"id"`
	Friendly_name   string  `json:"friendly_name"`
	Url             string  `json:"url"`
	Monitor_type    int     `json:"type"`
	Sub_type        *int    `json:"sub_type,omitempty"`
	Keyword_type    *int    `json:"keyword_type,omitempty"`
	Keyword_value   *string `json:"keyword_value,omitempty"`
	Http_username   *string `json:"http_username,omitempty"`
	Http_password   *string `json:"http_password,omitempty"`
	Port            *string `json:"http_port,omitempty"`
	Interval        *int    `json:"interval,omitempty"`
	Status          *int    `json:"status,omitempty"`
	Create_datetime *int    `json:"create_datetime,omitempty"`
	Monitor_group   *int    `json:"monitor_group,omitempty"`
	Is_group_main   *int    `json:"is_group_main,omitempty"`
	Logs            []Log   `json:"logs,omitempty"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type MonitorResp struct {
	Stat       string     `json:"stat"`
	Pagination Pagination `json:"pagination"`
	Monitors   []Monitor  `json:"monitors"`
}

type CreatedMonitor struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type CreateMonitorResp struct {
	Stat    string         `json:"stat"`
	Monitor CreatedMonitor `json:"monitor"`
}

func (c *Client) _makeReq(path string, data *url.Values) (*http.Request, error) {
	data.Set("api_key", c.Api_key)
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

func (c *Client) GetMonitors() ([]Monitor, error) {
	data := url.Values{}

	req, err := c._makeReq("/getMonitors", &data)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//bodybuffer, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s", bodybuffer)
	var monitors_resp MonitorResp
	err = json.NewDecoder(resp.Body).Decode(&monitors_resp)
	if err != nil {
		return nil, err
	}
	return monitors_resp.Monitors, err
}

func _optionalInt(data *url.Values, key string, value *int) {
	if value != nil {
		data.Set(key, fmt.Sprintf("%d", *value))
	}
}

func _optionalString(data *url.Values, key string, value *string) {
	if value != nil {
		data.Set(key, fmt.Sprintf("%s", *value))
	}
}

func (c *Client) CreateMonitor(m *Monitor) error {
	data := url.Values{}
	req, err := c._makeReq("/newMonitor", &data)
	if err != nil {
		return err
	}
	data.Set("friendly_name", m.Friendly_name)
	data.Set("url", m.Url)
	data.Set("type", fmt.Sprintf("%d", m.Monitor_type))
	_optionalInt(&data, "sub_type", m.Sub_type)
	_optionalInt(&data, "keyword_type", m.Keyword_type)
	_optionalString(&data, "keyword_value", m.Keyword_value)
	_optionalString(&data, "http_username", m.Http_username)
	_optionalString(&data, "http_password", m.Http_password)
	_optionalString(&data, "port", m.Port)
	_optionalInt(&data, "interval", m.Interval)
	_optionalInt(&data, "status", m.Status)
	_optionalInt(&data, "create_datetime", m.Create_datetime)
	_optionalInt(&data, "monitor_group", m.Monitor_group)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var monitor_create_resp CreateMonitorResp
	err = json.NewDecoder(resp.Body).Decode(&monitor_create_resp)
	if err != nil {
		return err
	}
	m.Id = monitor_create_resp.Monitor.Id
	return nil
}
