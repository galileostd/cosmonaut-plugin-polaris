package client

import (
    "context"
    "net/url"
	"fmt"

    "github.com/galileostd/cosmonaut-plugin-polaris/internal/types"
)

// GenericClient handles Generic Tables API operations.
type GenericClient struct {
    *PolarisClient
}

// NewGenericClient creates a new Generic client.
func NewGenericClient(cfg PolarisClientConfig) *GenericClient {
    return &GenericClient{PolarisClient: NewPolarisClient(cfg)}
}

// ─── Generic Tables ────────────────────────────────────────────────────────────

// ListGenericTables lists all generic tables in a namespace.
func (c *GenericClient) ListGenericTables(ctx context.Context, namespace []string, pageToken string, pageSize int) (*types.ListGenericTablesResponse, error) {
    params := url.Values{}
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "generic-tables")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.ListGenericTablesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// CreateGenericTable creates a new generic table.
func (c *GenericClient) CreateGenericTable(ctx context.Context, namespace []string, req types.CreateGenericTableRequest) (*types.GenericTable, error) {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "generic-tables")
    var resp types.GenericTable
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// LoadGenericTable loads a generic table.
func (c *GenericClient) LoadGenericTable(ctx context.Context, namespace []string, tableName string) (*types.LoadGenericTableResponse, error) {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "generic-tables", tableName)
    var resp types.LoadGenericTableResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// DropGenericTable deletes a generic table.
func (c *GenericClient) DropGenericTable(ctx context.Context, namespace []string, tableName string) error {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "generic-tables", tableName)
    return c.doRequest(ctx, "DELETE", url, nil, nil)
}