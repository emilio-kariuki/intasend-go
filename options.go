package intasend

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring the Client.
type Option func(*Client) error

// WithPublishableKey sets the publishable (public) API key.
// Keys starting with "ISPubKey_test" indicate the sandbox environment.
func WithPublishableKey(key string) Option {
	return func(c *Client) error {
		c.publishableKey = key
		return nil
	}
}

// WithSecretKey sets the secret API key for authenticated requests.
// Keys starting with "ISSecretKey_test" indicate the sandbox environment.
func WithSecretKey(key string) Option {
	return func(c *Client) error {
		c.secretKey = key
		return nil
	}
}

// WithBaseURL overrides the auto-detected base URL.
// Use this if you need to point to a custom API endpoint.
func WithBaseURL(url string) Option {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
// This is useful for testing or custom transport configuration.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

// WithTimeout sets the request timeout duration.
// Default is 30 seconds.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}

// WithRetry configures the retry behavior for failed requests.
// Default is 3 retries with 1 second initial wait (exponential backoff).
func WithRetry(maxRetries int, waitTime time.Duration) Option {
	return func(c *Client) error {
		c.maxRetries = maxRetries
		c.retryWait = waitTime
		return nil
	}
}

// WithDebug enables debug logging of requests and responses.
func WithDebug(debug bool) Option {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		c.userAgent = ua
		return nil
	}
}

// WithSandbox forces the client to use the sandbox environment.
func WithSandbox() Option {
	return func(c *Client) error {
		c.baseURL = SandboxBaseURL
		return nil
	}
}

// WithProduction forces the client to use the production environment.
func WithProduction() Option {
	return func(c *Client) error {
		c.baseURL = ProductionBaseURL
		return nil
	}
}
