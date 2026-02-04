package intasend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	headerAuthorization     = "Authorization"
	headerContentType       = "Content-Type"
	headerPublicAPIKey      = "X-IntaSend-Public-API-Key"
	headerIntaSendPublicKey = "INTASEND_PUBLIC_API_KEY"
	headerUserAgent         = "User-Agent"

	contentTypeJSON = "application/json"
)

// requestConfig holds configuration for a single request.
type requestConfig struct {
	method        string
	path          string
	body          interface{}
	result        interface{}
	requiresAuth  bool
	publicKeyOnly bool
}

// doRequest performs an HTTP request with retries and error handling.
func (c *Client) doRequest(ctx context.Context, cfg *requestConfig) error {
	var bodyBytes []byte
	var err error

	if cfg.body != nil {
		bodyBytes, err = json.Marshal(cfg.body)
		if err != nil {
			return fmt.Errorf("intasend: failed to marshal request body: %w", err)
		}
	}

	url := c.baseURL + cfg.path

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := c.retryWait * time.Duration(1<<(attempt-1))
			if c.debug {
				log.Printf("[IntaSend] Retry attempt %d after %v", attempt, waitTime)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, cfg.method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("intasend: failed to create request: %w", err)
		}

		req.Header.Set(headerContentType, contentTypeJSON)
		req.Header.Set(headerUserAgent, c.userAgent)

		if c.publishableKey != "" {
			req.Header.Set(headerPublicAPIKey, c.publishableKey)
			req.Header.Set(headerIntaSendPublicKey, c.publishableKey)
		}

		if cfg.requiresAuth && c.secretKey != "" {
			req.Header.Set(headerAuthorization, "Bearer "+c.secretKey)
		}

		if c.debug {
			log.Printf("[IntaSend] %s %s", cfg.method, url)
			if bodyBytes != nil {
				log.Printf("[IntaSend] Request Body: %s", string(bodyBytes))
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = &NetworkError{Err: err, Message: "request failed"}
			if c.debug {
				log.Printf("[IntaSend] Network error: %v", err)
			}
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = &NetworkError{Err: err, Message: "failed to read response"}
			if c.debug {
				log.Printf("[IntaSend] Failed to read response: %v", err)
			}
			continue
		}

		if c.debug {
			log.Printf("[IntaSend] Response Status: %d", resp.StatusCode)
			log.Printf("[IntaSend] Response Body: %s", string(respBody))
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			apiErr := &APIError{HTTPStatusCode: resp.StatusCode}
			if err := json.Unmarshal(respBody, apiErr); err != nil {
				apiErr.Message = string(respBody)
			}

			// Don't retry client errors (except rate limiting)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
				return apiErr
			}
			lastErr = apiErr
			continue
		}

		if cfg.result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, cfg.result); err != nil {
				return fmt.Errorf("intasend: failed to unmarshal response: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, &requestConfig{
		method:       http.MethodGet,
		path:         path,
		result:       result,
		requiresAuth: true,
	})
}

// post performs a POST request with authentication.
func (c *Client) post(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, &requestConfig{
		method:       http.MethodPost,
		path:         path,
		body:         body,
		result:       result,
		requiresAuth: true,
	})
}

// postPublic performs a POST request using only the public key (no auth).
func (c *Client) postPublic(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, &requestConfig{
		method:        http.MethodPost,
		path:          path,
		body:          body,
		result:        result,
		requiresAuth:  false,
		publicKeyOnly: true,
	})
}
