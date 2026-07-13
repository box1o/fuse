package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	mu      sync.RWMutex
	baseURL string
	http    *http.Client
}

type ResponseError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("server returned %s: %s", e.Status, e.Body)
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) SetBaseURL(baseURL string) error {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	parsed, err := url.Parse(baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return fmt.Errorf("invalid API URL %q", baseURL)
	}

	c.mu.Lock()
	c.baseURL = baseURL
	c.mu.Unlock()
	return nil
}

func (c *Client) endpoint(path string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL + path
}

func (c *Client) Do(ctx context.Context, method, path, token string, requestBody, responseBody any) (int, error) {
	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return 0, err
		}
		body = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.endpoint(path), body)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Accept", "application/json")
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := c.http.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return response.StatusCode, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return response.StatusCode, &ResponseError{StatusCode: response.StatusCode, Status: response.Status, Body: string(payload)}
	}
	if responseBody != nil && len(payload) > 0 {
		if err := json.Unmarshal(payload, responseBody); err != nil {
			return response.StatusCode, err
		}
	}
	return response.StatusCode, nil
}
