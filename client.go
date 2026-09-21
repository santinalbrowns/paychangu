package paychangu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.paychangu.com"

// PayChangu is the API client.
type PayChangu struct {
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

// Option configures a PayChangu client.
type Option func(*PayChangu)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *PayChangu) {
		if c != nil {
			p.httpClient = c
		}
	}
}

// WithBaseURL overrides the API base URL (useful for tests).
func WithBaseURL(base string) Option {
	return func(p *PayChangu) {
		if base != "" {
			p.baseURL = strings.TrimRight(base, "/")
		}
	}
}

// New creates a PayChangu client authenticated with the given secret key.
func New(secretKey string, opts ...Option) *PayChangu {
	p := &PayChangu{
		secretKey: secretKey,
		baseURL:   defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func isSuccessStatus(status string) bool {
	s := strings.ToLower(strings.TrimSpace(status))
	return s == "success" || s == "successful"
}

func (p *PayChangu) do(ctx context.Context, method, path string, body any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}

	fullURL := p.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, raw)
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func parseAPIError(statusCode int, raw []byte) error {
	var envelope struct {
		Status  string          `json:"status"`
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Message) > 0 {
		var msg string
		if json.Unmarshal(envelope.Message, &msg) == nil && msg != "" {
			return fmt.Errorf("API error (%d): %s", statusCode, msg)
		}
		var validation map[string][]string
		if json.Unmarshal(envelope.Message, &validation) == nil && len(validation) > 0 {
			parts := make([]string, 0)
			for field, msgs := range validation {
				for _, m := range msgs {
					parts = append(parts, fmt.Sprintf("%s: %s", field, m))
				}
			}
			return fmt.Errorf("API error (%d): %s", statusCode, strings.Join(parts, "; "))
		}
	}

	if len(raw) == 0 {
		return fmt.Errorf("API request failed with status %d", statusCode)
	}
	return fmt.Errorf("API request failed with status %d: %s", statusCode, string(raw))
}

func (p *PayChangu) get(ctx context.Context, path string, out any) error {
	return p.do(ctx, http.MethodGet, path, nil, out)
}

func (p *PayChangu) post(ctx context.Context, path string, body, out any) error {
	return p.do(ctx, http.MethodPost, path, body, out)
}

func (p *PayChangu) put(ctx context.Context, path string, body, out any) error {
	return p.do(ctx, http.MethodPut, path, body, out)
}

func (p *PayChangu) delete(ctx context.Context, path string, out any) error {
	return p.do(ctx, http.MethodDelete, path, nil, out)
}

func queryPath(path string, values url.Values) string {
	if len(values) == 0 {
		return path
	}
	return path + "?" + values.Encode()
}
