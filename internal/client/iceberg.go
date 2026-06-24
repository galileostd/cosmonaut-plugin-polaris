package client

import (
    "context"
    "fmt"
    "net/url"
	"net/http"

    "github.com/galileostd/cosmonaut-plugin-polaris/internal/types"
)

// IcebergClient handles Iceberg REST Catalog API operations.
type IcebergClient struct {
    *PolarisClient
}

// NewIcebergClient creates a new Iceberg REST client.
func NewIcebergClient(cfg PolarisClientConfig) *IcebergClient {
    return &IcebergClient{PolarisClient: NewPolarisClient(cfg)}
}

// ─── Namespaces ────────────────────────────────────────────────────────────────

// ListNamespaces lists all namespaces, optionally with a parent.
func (c *IcebergClient) ListNamespaces(ctx context.Context, parent []string, pageToken string, pageSize int) (*types.ListNamespacesResponse, error) {
    params := url.Values{}
    if parent != nil && len(parent) > 0 {
        params.Set("parent", encodeNamespace(parent))
    }
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.IcebergURL("namespaces")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.ListNamespacesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// CreateNamespace creates a new namespace.
func (c *IcebergClient) CreateNamespace(ctx context.Context, namespace []string, properties map[string]string) (*types.CreateNamespaceResponse, error) {
    req := types.CreateNamespaceRequest{
        Namespace:  namespace,
        Properties: properties,
    }
    url := c.IcebergURL("namespaces")
    var resp types.CreateNamespaceResponse
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// LoadNamespace loads a namespace's metadata.
func (c *IcebergClient) LoadNamespace(ctx context.Context, namespace []string) (*types.GetNamespaceResponse, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace))
    var resp types.GetNamespaceResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// NamespaceExists checks if a namespace exists.
func (c *IcebergClient) NamespaceExists(ctx context.Context, namespace []string) (bool, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace))
    req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
    if err != nil {
        return false, err
    }
    req.Header.Set("Authorization", "Bearer "+c.token)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    return resp.StatusCode == 204, nil
}

// DropNamespace deletes a namespace (must be empty).
func (c *IcebergClient) DropNamespace(ctx context.Context, namespace []string) error {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace))
    return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// UpdateNamespaceProperties updates or removes properties on a namespace.
func (c *IcebergClient) UpdateNamespaceProperties(ctx context.Context, namespace []string, removals []string, updates map[string]string) (*types.UpdateNamespacePropertiesResponse, error) {
    req := types.UpdateNamespacePropertiesRequest{
        Removals: removals,
        Updates:  updates,
    }
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "properties")
    var resp types.UpdateNamespacePropertiesResponse
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// ─── Tables ────────────────────────────────────────────────────────────────────

// ListTables lists all tables in a namespace.
func (c *IcebergClient) ListTables(ctx context.Context, namespace []string, pageToken string, pageSize int) (*types.ListTablesResponse, error) {
    params := url.Values{}
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.ListTablesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// CreateTable creates a new table.
func (c *IcebergClient) CreateTable(ctx context.Context, namespace []string, req types.CreateTableRequest) (*types.LoadTableResult, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables")
    var resp types.LoadTableResult
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// LoadTable loads a table's metadata.
func (c *IcebergClient) LoadTable(ctx context.Context, namespace []string, tableName string, snapshots string) (*types.LoadTableResult, error) {
    params := url.Values{}
    if snapshots != "" {
        params.Set("snapshots", snapshots)
    }

    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName)
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.LoadTableResult
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// TableExists checks if a table exists.
func (c *IcebergClient) TableExists(ctx context.Context, namespace []string, tableName string) (bool, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName)
    req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
    if err != nil {
        return false, err
    }
    req.Header.Set("Authorization", "Bearer "+c.token)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    return resp.StatusCode == 204, nil
}

// UpdateTable commits updates to a table.
func (c *IcebergClient) UpdateTable(ctx context.Context, namespace []string, tableName string, req types.CommitTableRequest) (*types.CommitTableResponse, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName)
    var resp types.CommitTableResponse
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// DropTable deletes a table.
func (c *IcebergClient) DropTable(ctx context.Context, namespace []string, tableName string, purgeRequested bool) error {
    params := url.Values{}
    if purgeRequested {
        params.Set("purgeRequested", "true")
    }

    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName)
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// RenameTable renames a table (can move across namespaces).
func (c *IcebergClient) RenameTable(ctx context.Context, source, destination types.TableIdentifier) error {
    req := types.RenameTableRequest{
        Source:      source,
        Destination: destination,
    }
    url := c.IcebergURL("tables", "rename")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// RegisterTable registers a table from a metadata location.
func (c *IcebergClient) RegisterTable(ctx context.Context, namespace []string, req types.RegisterTableRequest) (*types.LoadTableResult, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "register")
    var resp types.LoadTableResult
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// LoadCredentials loads vended credentials for a table.
func (c *IcebergClient) LoadCredentials(ctx context.Context, namespace []string, tableName string) (*types.LoadCredentialsResponse, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName, "credentials")
    var resp types.LoadCredentialsResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// ReportMetrics sends a metrics report.
func (c *IcebergClient) ReportMetrics(ctx context.Context, namespace []string, tableName string, req types.ReportMetricsRequest) error {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName, "metrics")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// CommitTransaction commits multiple table updates atomically.
func (c *IcebergClient) CommitTransaction(ctx context.Context, req types.CommitTransactionRequest) error {
    url := c.IcebergURL("transactions", "commit")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// ─── Views ────────────────────────────────────────────────────────────────────

// ListViews lists all views in a namespace.
func (c *IcebergClient) ListViews(ctx context.Context, namespace []string, pageToken string, pageSize int) (*types.ListTablesResponse, error) {
    params := url.Values{}
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.ListTablesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// CreateView creates a new view.
func (c *IcebergClient) CreateView(ctx context.Context, namespace []string, req types.CreateViewRequest) (*types.LoadViewResult, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views")
    var resp types.LoadViewResult
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// LoadView loads a view's metadata.
func (c *IcebergClient) LoadView(ctx context.Context, namespace []string, viewName string) (*types.LoadViewResult, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views", viewName)
    var resp types.LoadViewResult
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// ViewExists checks if a view exists.
func (c *IcebergClient) ViewExists(ctx context.Context, namespace []string, viewName string) (bool, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views", viewName)
    req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
    if err != nil {
        return false, err
    }
    req.Header.Set("Authorization", "Bearer "+c.token)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    return resp.StatusCode == 204, nil
}

// ReplaceView replaces a view (commit updates).
func (c *IcebergClient) ReplaceView(ctx context.Context, namespace []string, viewName string, req types.CommitViewRequest) (*types.LoadViewResult, error) {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views", viewName)
    var resp types.LoadViewResult
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// DropView deletes a view.
func (c *IcebergClient) DropView(ctx context.Context, namespace []string, viewName string) error {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "views", viewName)
    return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// RenameView renames a view.
func (c *IcebergClient) RenameView(ctx context.Context, source, destination types.TableIdentifier) error {
    req := types.RenameTableRequest{
        Source:      source,
        Destination: destination,
    }
    url := c.IcebergURL("views", "rename")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// ─── Notifications ────────────────────────────────────────────────────────────

// SendNotification sends a notification to a table.
func (c *IcebergClient) SendNotification(ctx context.Context, namespace []string, tableName string, req types.NotificationRequest) error {
    url := c.IcebergURL("namespaces", encodeNamespace(namespace), "tables", tableName, "notifications")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func encodeNamespace(parts []string) string {
    if len(parts) == 0 {
        return ""
    }
    result := parts[0]
    for i := 1; i < len(parts); i++ {
        result += "\x1F" + parts[i]
    }
    return result
}

// GetConfig calls GET /api/catalog/v1/config — used for health checks.
func (c *IcebergClient) GetConfig(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/catalog/v1/config?warehouse=%s", c.endpoint, c.catalog)
	var result map[string]interface{}
	if err := c.doRequest(ctx, "GET", url, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
