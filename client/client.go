package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	token      string
}

type InstallationResponse struct {
	ID int64 `json:"id"`
}

type AccessTokenRequest struct {
	Repositories []string          `json:"repositories,omitempty"`
	Permissions  map[string]string `json:"permissions,omitempty"`
}

type AccessTokenResponse struct {
	Token string `json:"token"`
}

func New(baseURL, jwt string) *Client {
	return &Client{
		BaseURL: baseURL,
		token:   jwt,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) GetInstallationByOwner(owner string) (*http.Response, error) {
	path := fmt.Sprintf("users/%s/installation", owner)
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

func (c *Client) GetInstallationAccessToken(installationID int64, permissions map[string]string, repos []string) (*AccessTokenResponse, error) {
	path := fmt.Sprintf("app/installations/%d/access_tokens", installationID)

	payload := AccessTokenRequest{
		Repositories: repos,
		Permissions:  permissions,
	}

	req, err := c.newRequest(http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get token: %s, body: %s", resp.Status, string(body))
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (c *Client) DeleteInstallationAccessToken() error {
	path := "installation/token"

	req, err := c.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete token: %s", resp.Status)
	}

	return nil
}
