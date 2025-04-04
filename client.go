package xui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/novikoff-vvs/xui/dto"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client represents the 3X-UI API client
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	sessionID  string
}

// NewClient creates a new 3X-UI API client
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		username:   username,
		password:   password,
		httpClient: &http.Client{},
	}
}

// Login authenticates with the 3X-UI panel and stores the session cookie
func (c *Client) Login() error {
	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", c.baseURL), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %s", resp.Status)
	}

	// Parse response
	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("login failed: %s", result.Msg)
	}

	// Extract session cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "3x-ui" {
			c.sessionID = cookie.Value
			break
		}
	}

	if c.sessionID == "" {
		return fmt.Errorf("session cookie not found in login response")
	}

	return nil
}

// Inbound represents an inbound configuration

// GetInbounds retrieves a list of all inbounds
func (c *Client) GetInbounds() ([]dto.Inbound, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/panel/api/inbounds/list", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool          `json:"success"`
		Msg     string        `json:"msg"`
		Obj     []dto.Inbound `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get inbounds: %s", result.Msg)
	}

	return result.Obj, nil
}

// GetInbound retrieves details for a specific inbound
func (c *Client) GetInbound(inboundID int) (*dto.Inbound, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/panel/api/inbounds/get/%d", c.baseURL, inboundID), nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool        `json:"success"`
		Msg     string      `json:"msg"`
		Obj     dto.Inbound `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get inbound: %s", result.Msg)
	}

	return &result.Obj, nil
}

// GetClientTraffics retrieves traffic information for a client by email
func (c *Client) GetClientTraffics(email string) (*dto.ClientTraffic, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/panel/api/inbounds/getClientTraffics/%s", c.baseURL, email), nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool               `json:"success"`
		Msg     string             `json:"msg"`
		Obj     *dto.ClientTraffic `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get client traffics: %s", result.Msg)
	}

	return result.Obj, nil
}

// GetClientTrafficsByID retrieves traffic information for a client by UUID
func (c *Client) GetClientTrafficsByID(uuid string) ([]dto.ClientTraffic, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/panel/api/inbounds/getClientTrafficsById/%s", c.baseURL, uuid), nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool                `json:"success"`
		Msg     string              `json:"msg"`
		Obj     []dto.ClientTraffic `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get client traffics by ID: %s", result.Msg)
	}

	return result.Obj, nil
}

// CreateBackup triggers creation of a system backup
func (c *Client) CreateBackup() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/panel/api/inbounds/createbackup", c.baseURL), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to create backup: %s", result.Msg)
	}

	return nil
}

// GetClientIPs retrieves IP records for a client
func (c *Client) GetClientIPs(email string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/clientIps/%s", c.baseURL, email), nil)
	if err != nil {
		return "", err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Obj     string `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf("failed to get client IPs: %s", result.Msg)
	}

	return result.Obj, nil
}

// AddInbound adds a new inbound configuration
func (c *Client) AddInbound(inbound dto.Inbound) (*dto.Inbound, error) {
	body, err := json.Marshal(inbound)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/add", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)
	req.Header.Add("Content-Type", "application/json")

	var result struct {
		Success bool        `json:"success"`
		Msg     string      `json:"msg"`
		Obj     dto.Inbound `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to add inbound: %s", result.Msg)
	}

	return &result.Obj, nil
}

// AddClientToInbound adds a new client to an existing inbound
func (c *Client) AddClientToInbound(inboundID int, clientSettings string) error {
	payload := struct {
		ID       int    `json:"id"`
		Settings string `json:"settings"`
	}{
		ID:       inboundID,
		Settings: clientSettings,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/addClient", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.setAuthHeader(req)
	req.Header.Add("Content-Type", "application/json")

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to add client: %s", result.Msg)
	}

	return nil
}

// UpdateInbound updates an existing inbound configuration
func (c *Client) UpdateInbound(inboundID int, inbound dto.Inbound) (*dto.Inbound, error) {
	body, err := json.Marshal(inbound)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/update/%d", c.baseURL, inboundID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)
	req.Header.Add("Content-Type", "application/json")

	var result struct {
		Success bool        `json:"success"`
		Msg     string      `json:"msg"`
		Obj     dto.Inbound `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to update inbound: %s", result.Msg)
	}

	return &result.Obj, nil
}

// UpdateClient updates an existing client configuration
func (c *Client) UpdateClient(uuid string, inboundID int, clientSettings string) error {
	payload := struct {
		ID       int    `json:"id"`
		Settings string `json:"settings"`
	}{
		ID:       inboundID,
		Settings: clientSettings,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/updateClient/%s", c.baseURL, uuid), bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.setAuthHeader(req)
	req.Header.Add("Content-Type", "application/json")

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to update client: %s", result.Msg)
	}

	return nil
}

// ClearClientIPs clears IP records for a client
func (c *Client) ClearClientIPs(email string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/clearClientIps/%s", c.baseURL, email), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to clear client IPs: %s", result.Msg)
	}

	return nil
}

// ResetAllTraffics resets traffic statistics for all inbounds
func (c *Client) ResetAllTraffics() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/resetAllTraffics", c.baseURL), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to reset all traffics: %s", result.Msg)
	}

	return nil
}

// ResetAllClientTraffics resets traffic statistics for all clients in an inbound
func (c *Client) ResetAllClientTraffics(inboundID int) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/resetAllClientTraffics/%d", c.baseURL, inboundID), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to reset all client traffics: %s", result.Msg)
	}

	return nil
}

// ResetClientTraffic resets traffic statistics for a specific client
func (c *Client) ResetClientTraffic(inboundID int, email string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/%d/resetClientTraffic/%s", c.baseURL, inboundID, email), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to reset client traffic: %s", result.Msg)
	}

	return nil
}

// DeleteClient deletes a client from an inbound
func (c *Client) DeleteClient(inboundID int, uuid string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/%d/delClient/%s", c.baseURL, inboundID, uuid), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to delete client: %s", result.Msg)
	}

	return nil
}

// DeleteInbound deletes an inbound configuration
func (c *Client) DeleteInbound(inboundID int) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/del/%d", c.baseURL, inboundID), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Obj     int    `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to delete inbound: %s", result.Msg)
	}

	return nil
}

// DeleteDepletedClients deletes all depleted clients from an inbound
func (c *Client) DeleteDepletedClients(inboundID int) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/delDepletedClients/%d", c.baseURL, inboundID), nil)
	if err != nil {
		return err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to delete depleted clients: %s", result.Msg)
	}

	return nil
}

// GetOnlineClients retrieves a list of online clients
func (c *Client) GetOnlineClients() ([]string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/panel/api/inbounds/onlines", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	c.setAuthHeader(req)

	var result struct {
		Success bool     `json:"success"`
		Msg     string   `json:"msg"`
		Obj     []string `json:"obj"`
	}

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get online clients: %s", result.Msg)
	}

	return result.Obj, nil
}

// setAuthHeader sets the session cookie for authentication
func (c *Client) setAuthHeader(req *http.Request) {
	if c.sessionID != "" {
		req.AddCookie(&http.Cookie{
			Name:  "3x-ui",
			Value: c.sessionID,
		})
	}
}

// doRequest performs the HTTP request and decodes the response
func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
