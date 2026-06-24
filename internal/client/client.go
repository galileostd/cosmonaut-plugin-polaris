package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

const (
    DefaultTimeout = 30 * time.Second
)

// PolarisClient is the base HTTP client for all Polaris APIs.
type PolarisClient struct {
    endpoint   string
    httpClient *http.Client
    token      string
    catalog    string
    prefix     string
}

// PolarisClientConfig holds configuration for the Polaris client.
type PolarisClientConfig struct {
    Endpoint string
    Token    string
    Catalog  string // Default catalog name
    Prefix   string // API prefix, defaults to "v1"
    Timeout  time.Duration
}

// NewPolarisClient creates a new Polaris client.
func NewPolarisClient(cfg PolarisClientConfig) *PolarisClient {
    if cfg.Prefix == "" {
        cfg.Prefix = "v1"
    }
    if cfg.Timeout == 0 {
        cfg.Timeout = DefaultTimeout
    }
    return &PolarisClient{
        endpoint:   strings.TrimRight(cfg.Endpoint, "/"),
        httpClient: &http.Client{Timeout: cfg.Timeout},
        token:      cfg.Token,
        catalog:    cfg.Catalog,
        prefix:     cfg.Prefix,
    }
}

// SetToken updates the bearer token.
func (c *PolarisClient) SetToken(token string) {
    c.token = token
}

// doRequest performs an HTTP request and decodes the response.
func (c *PolarisClient) doRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
    var reqBody io.Reader
    if body != nil {
        b, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("marshaling request body: %w", err)
        }
        reqBody = bytes.NewReader(b)
    }

    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return fmt.Errorf("creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("executing request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
    }

    if result != nil {
        if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
            return fmt.Errorf("decoding response: %w", err)
        }
    }

    return nil
}

// IcebergURL builds a URL for the Iceberg REST API.
func (c *PolarisClient) IcebergURL(path string, parts ...string) string {
    base := fmt.Sprintf("%s/api/catalog/%s/%s", c.endpoint, c.prefix, c.catalog)
    if path != "" {
        base = base + "/" + strings.TrimLeft(path, "/")
    }
    for _, p := range parts {
        base = base + "/" + p
    }
    return base
}

// ManagementURL builds a URL for the Polaris Management API.
func (c *PolarisClient) ManagementURL(path string, parts ...string) string {
    base := fmt.Sprintf("%s/api/management/v1/%s", c.endpoint, strings.TrimLeft(path, "/"))
    for _, p := range parts {
        base = base + "/" + p
    }
    return base
}

// PolarisURL builds a URL for Polaris-specific APIs (policies, generic tables).
func (c *PolarisClient) PolarisURL(path string, parts ...string) string {
    base := fmt.Sprintf("%s/api/catalog/polaris/%s/%s/%s", c.endpoint, c.prefix, c.catalog, strings.TrimLeft(path, "/"))
    for _, p := range parts {
        base = base + "/" + p
    }
    return base
}
