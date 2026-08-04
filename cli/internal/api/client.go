package api

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

const maxResponseSize = 2 << 20

type Client struct {
	mu      sync.RWMutex
	baseURL string
	http    *http.Client
}

type ResponseError struct {
	StatusCode int
	Status     string
	Code       string
	Message    string
}

func (e *ResponseError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("server returned %s", e.Status)
	}

	return fmt.Sprintf("server returned %s: %s", e.Status, e.Message)
}

func NewClient(baseURL string) (*Client, error) {
	client := &Client{http: &http.Client{Timeout: 30 * time.Second}}
	if err := client.SetBaseURL(baseURL); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) SetBaseURL(baseURL string) error {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("invalid Fuse API URL %q", baseURL)
	}

	c.mu.Lock()
	c.baseURL = baseURL
	c.mu.Unlock()
	return nil
}

func (c *Client) Do(ctx context.Context, method, path, token string, requestBody, responseBody any) (int, error) {
	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return 0, fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(ctx, method, c.endpoint(path), body)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
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
		return 0, fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize))
	if err != nil {
		return response.StatusCode, fmt.Errorf("read response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return response.StatusCode, decodeResponseError(response, payload)
	}
	if responseBody != nil && len(payload) > 0 {
		if err := json.Unmarshal(payload, responseBody); err != nil {
			return response.StatusCode, fmt.Errorf("decode response: %w", err)
		}
	}

	return response.StatusCode, nil
}

func (c *Client) endpoint(path string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL + "/" + strings.TrimLeft(path, "/")
}

func decodeResponseError(response *http.Response, payload []byte) error {
	var body struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}
	_ = json.Unmarshal(payload, &body)

	message := body.Detail
	if message == "" {
		message = body.Message
	}
	if message == "" {
		message = body.Error
	}
	if message == "" {
		message = strings.TrimSpace(string(payload))
	}

	return &ResponseError{StatusCode: response.StatusCode, Status: response.Status, Code: body.Error, Message: message}
}
