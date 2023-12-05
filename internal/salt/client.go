// SPDX-License-Identifier: Apache-2.0

package salt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type loginResponse struct {
	Return []struct {
		Token string `json:"token"`
	} `json:"return"`
}

type Client struct {
	endpoint string
	client   *http.Client
	token    string
}

func New(endpoint string, username string, password string) (*Client, error) {
	token, err := authenticate(endpoint, username, password)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: endpoint,
		client:   http.DefaultClient,
		token:    token,
	}, nil
}

func authenticate(endpoint string, username string, password string) (string, error) {
	c := http.DefaultClient

	data := url.Values{}
	data.Add("eauth", "pam")
	data.Add("username", username)
	data.Add("password", password)

	request, err := http.NewRequest(http.MethodPost, endpoint+"/login", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("NewRequest: %w", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("User-Agent", "terraform-provider-salt/0.0.1")
	resp, err := c.Do(request)
	if err != nil {
		return "", fmt.Errorf("c.Do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not OK: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll: %w", err)
	}

	login := &loginResponse{}
	err = json.Unmarshal(body, login)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return login.Return[0].Token, nil
}

func (s Client) Do(req *http.Request) (*http.Response, error) {
	// TODO get version from provider
	req.Header.Add("User-Agent", "terraform-provider-salt/0.0.1")
	req.Header.Add("X-Auth-Token", s.token)

	return s.client.Do(req)
}
