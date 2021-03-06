/*
Copyright 2018 Clint Byrum

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.package uptimerobot
*/

package uptimerobot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

var MonitorTypeToName = map[int]string{
	MONITOR_TYPE_HTTP:    "HTTP(s)",
	MONITOR_TYPE_KEYWORD: "keyword",
	MONITOR_TYPE_PING:    "ping",
	MONITOR_TYPE_PORT:    "port",
}

var MonitorTypeNames = map[string]int{
	"HTTP":    MONITOR_TYPE_HTTP,
	"HTTPS":   MONITOR_TYPE_HTTP,
	"HTTP(s)": MONITOR_TYPE_HTTP,
	"http":    MONITOR_TYPE_HTTP,
	"https":   MONITOR_TYPE_HTTP,
	"keyword": MONITOR_TYPE_KEYWORD,
	"ping":    MONITOR_TYPE_PING,
	"port":    MONITOR_TYPE_PORT,
}

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
	Sub_type        *string `json:"sub_type,omitempty"`
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

type ChangedMonitor struct {
	Id     int  `json:"id"`
	Status *int `json:"status,omitempty"`
}

type ChangeMonitorResp struct {
	Stat    string         `json:"stat"`
	Monitor ChangedMonitor `json:"monitor"`
}

func (c *Client) makeReq(path string, data *url.Values) (*http.Request, error) {
	data.Set("api_key", c.Api_key)
	data.Set("format", "json")
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	log.Print("[TRACE] Sending req to ", u)
	body := data.Encode()
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *Client) GetMonitors(ids []int) ([]Monitor, error) {
	data := url.Values{}
	if len(ids) > 0 {
		var str_ids []string
		for i := 0; i < len(ids); i++ {
			id := ids[i]
			str_ids = append(str_ids, fmt.Sprintf("%d", id))
		}
		data.Set("monitors", strings.Join(str_ids, "-"))
	}

	req, err := c.makeReq("getMonitors", &data)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	log.Print("[TRACE] ", buf.String())
	var monitors_resp MonitorResp
	err = json.NewDecoder(buf).Decode(&monitors_resp)
	if err != nil {
		return nil, err
	}
	return monitors_resp.Monitors, err
}

func optionalInt(data *url.Values, key string, value *int) {
	if value != nil {
		data.Set(key, fmt.Sprintf("%d", *value))
	}
}

func optionalString(data *url.Values, key string, value *string) {
	if value != nil {
		data.Set(key, fmt.Sprintf("%s", *value))
	}
}

func (c *Client) setCommonData(data *url.Values, m *Monitor) {
	data.Set("friendly_name", m.Friendly_name)
	data.Set("url", m.Url)
	optionalString(data, "sub_type", m.Sub_type)
	optionalInt(data, "keyword_type", m.Keyword_type)
	optionalString(data, "keyword_value", m.Keyword_value)
	optionalString(data, "http_username", m.Http_username)
	optionalString(data, "http_password", m.Http_password)
	optionalString(data, "port", m.Port)
	optionalInt(data, "interval", m.Interval)
	optionalInt(data, "status", m.Status)
}

func (c *Client) CreateMonitor(m *Monitor) error {
	data := url.Values{}
	req, err := c.makeReq("newMonitor", &data)
	if err != nil {
		log.Print("[DEBUG] error during makeReq = ", err)
		return err
	}
	// Only allowed to set this on create
	data.Set("type", fmt.Sprintf("%d", m.Monitor_type))
	change_resp, err := c.changeMonitor(req, &data, m)
	log.Print("[TRACE] resp = ", change_resp)
	if err != nil {
		log.Print("[DEBUG] error during changeMonitor = ", err)
		return err
	}
	m.Id = change_resp.Monitor.Id
	return nil
}

func (c *Client) EditMonitor(m *Monitor) error {
	if m.Id == 0 {
		return errors.New("Id is required to edit")
	}
	data := url.Values{}
	req, err := c.makeReq("editMonitor", &data)
	if err != nil {
		return err
	}
	_, err = c.changeMonitor(req, &data, m)
	return err
}

func (c *Client) DeleteMonitor(id int) error {
	data := url.Values{}
	req, err := c.makeReq("deleteMonitor", &data)
	if err != nil {
		return err
	}
	data.Set("id", fmt.Sprintf("%d", id))
	_, err = c.changeMonitor(req, &data, nil)
	return err
}

func (c *Client) changeMonitor(req *http.Request, data *url.Values, m *Monitor) (ChangeMonitorResp, error) {
	var monitor_change_resp ChangeMonitorResp
	if m != nil {
		c.setCommonData(data, m)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Print("[DEBUG] error during Httpclient.Do: ", err)
		return monitor_change_resp, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Print("[ERROR] got non-200 from server: ", resp.Status)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return monitor_change_resp, errors.New(buf.String())
	}

	err = json.NewDecoder(resp.Body).Decode(&monitor_change_resp)
	if err != nil {
		return monitor_change_resp, err
	}
	if monitor_change_resp.Stat != "ok" {
		log.Print("[DEBUG] Stat = ", monitor_change_resp.Stat)
		return monitor_change_resp, errors.New("Server reports request not ok!")
	}
	return monitor_change_resp, err
}
