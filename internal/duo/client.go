package duo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
)

type Client struct {
	duoapi.DuoApi
}

func NewClient(ikey, skey, host string) *Client {
	return &Client{
		*duoapi.NewDuoApi(ikey, skey, host, "terraform-provider-duo"),
	}
}

func (c *Client) DoRequest(method, path string, params url.Values, target interface{}) error {
	var resp *http.Response
	var err error

	resp, _, err = c.SignedCall(method, path, params)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var duoErr struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Stat    string `json:"stat"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&duoErr); err == nil && duoErr.Stat == "FAIL" {
			return fmt.Errorf("Duo API error %d: %s", duoErr.Code, duoErr.Message)
		}
		return fmt.Errorf("Duo API returned status %d", resp.StatusCode)
	}

	var wrapper struct {
		Response interface{} `json:"response"`
		Stat     string      `json:"stat"`
	}
	wrapper.Response = target

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return err
	}

	if wrapper.Stat != "OK" {
		return fmt.Errorf("Duo API returned non-OK stat: %s", wrapper.Stat)
	}

	return nil
}

// Group Types
type Group struct {
	GroupId     string `json:"group_id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Status      string `json:"status"`
	PushEnabled bool   `json:"push_enabled"`
	SmsEnabled  bool   `json:"sms_enabled"`
}

func (c *Client) CreateGroup(params url.Values) (*Group, error) {
	var group Group
	err := c.DoRequest("POST", "/admin/v1/groups", params, &group)
	return &group, err
}

func (c *Client) GetGroups() ([]Group, error) {
	var groups []Group
	err := c.DoRequest("GET", "/admin/v1/groups", nil, &groups)
	return groups, err
}

func (c *Client) UpdateGroup(id string, params url.Values) (*Group, error) {
	var group Group
	err := c.DoRequest("POST", "/admin/v1/groups/"+id, params, &group)
	return &group, err
}

func (c *Client) DeleteGroup(id string) error {
	return c.DoRequest("DELETE", "/admin/v1/groups/"+id, nil, nil)
}

// Application Types
type Application struct {
	IntegrationKey string `json:"integration_key"`
	Name           string `json:"name"`
	Type           string `json:"type"`
}

func (c *Client) CreateApplication(params url.Values) (*Application, error) {
	var app Application
	err := c.DoRequest("POST", "/admin/v1/apps", params, &app)
	return &app, err
}

func (c *Client) GetApplications() ([]Application, error) {
	var apps []Application
	err := c.DoRequest("GET", "/admin/v1/apps", nil, &apps)
	return apps, err
}

func (c *Client) UpdateApplication(id string, params url.Values) (*Application, error) {
	var app Application
	err := c.DoRequest("POST", "/admin/v1/apps/"+id, params, &app)
	return &app, err
}

func (c *Client) DeleteApplication(id string) error {
	return c.DoRequest("DELETE", "/admin/v1/apps/"+id, nil, nil)
}

// Policy Types
type Policy struct {
	PolicyKey string `json:"policy_key"`
	Name      string `json:"name"`
}

func (c *Client) CreatePolicy(params url.Values) (*Policy, error) {
	var policy Policy
	err := c.DoRequest("POST", "/admin/v1/policies", params, &policy)
	return &policy, err
}

func (c *Client) GetPolicies() ([]Policy, error) {
	var policies []Policy
	err := c.DoRequest("GET", "/admin/v1/policies", nil, &policies)
	return policies, err
}

func (c *Client) UpdatePolicy(id string, params url.Values) (*Policy, error) {
	var policy Policy
	err := c.DoRequest("POST", "/admin/v1/policies/"+id, params, &policy)
	return &policy, err
}

func (c *Client) DeletePolicy(id string) error {
	return c.DoRequest("DELETE", "/admin/v1/policies/"+id, nil, nil)
}

// User Types
type User struct {
	UserId    string  `json:"user_id"`
	Username  string  `json:"username"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	RealName  string  `json:"realname"`
	Email     string  `json:"email"`
	Status    string  `json:"status"`
	Notes     string  `json:"notes"`
	Alias1    string  `json:"alias1"`
	Alias2    string  `json:"alias2"`
	Alias3    string  `json:"alias3"`
	Alias4    string  `json:"alias4"`
	Groups    []Group `json:"groups"`
	Phones    []Phone `json:"phones"`
}

// Phone Types
type Phone struct {
	PhoneId   string `json:"phone_id"`
	Number    string `json:"number"`
	Extension string `json:"extension"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Platform  string `json:"platform"`
}

func (c *Client) CreateUser(params url.Values) (*User, error) {
	var user User
	err := c.DoRequest("POST", "/admin/v1/users", params, &user)
	return &user, err
}

func (c *Client) GetUser(id string) (*User, error) {
	var user User
	err := c.DoRequest("GET", "/admin/v1/users/"+id, nil, &user)
	return &user, err
}

func (c *Client) UpdateUser(id string, params url.Values) (*User, error) {
	var user User
	err := c.DoRequest("POST", "/admin/v1/users/"+id, params, &user)
	return &user, err
}

func (c *Client) DeleteUser(id string) error {
	return c.DoRequest("DELETE", "/admin/v1/users/"+id, nil, nil)
}

// Settings
type Settings struct {
	LockoutThreshold   int  `json:"lockout_threshold"`
	LockoutDuration    int  `json:"lockout_duration"`
	InactiveExpiration int  `json:"inactive_expiration"`
	UserApprove        bool `json:"user_approve"`
	UserTelephony      bool `json:"user_telephony"`
}

func (c *Client) GetSettings() (*Settings, error) {
	var s Settings
	err := c.DoRequest("GET", "/admin/v1/settings", nil, &s)
	return &s, err
}

func (c *Client) UpdateSettings(params url.Values) (*Settings, error) {
	var s Settings
	err := c.DoRequest("POST", "/admin/v1/settings", params, &s)
	return &s, err
}
