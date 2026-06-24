package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	pluginv1 "github.com/galileostd/cosmonaut-sdk/go/plugin/v1"

	"github.com/galileostd/cosmonaut-plugin-polaris/internal/client"
	"github.com/galileostd/cosmonaut-plugin-polaris/internal/types"
)

// ─── NAMESPACE HANDLERS ──────────────────────────────────────────────────────

func (p *Plugin) handleListNamespaces(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	parent, _ := req.Payload["parent"]
	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	var parentParts []string
	if parent != "" {
		parentParts = []string{parent}
	}

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.ListNamespaces(ctx, parentParts, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d namespaces", len(resp.Namespaces)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateNamespace(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var props map[string]string
	if propsJSON, ok := req.Payload["properties"]; ok {
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid properties format: %v", err),
			}, nil
		}
	}

	resp, err := c.CreateNamespace(ctx, namespace, props)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Namespace %v created", namespace),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadNamespace(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	resp, err := c.LoadNamespace(ctx, namespace)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded namespace %v", namespace),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleNamespaceExists(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	exists, err := c.NamespaceExists(ctx, namespace)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Namespace exists: %v", exists),
		Result:  map[string]string{"exists": fmt.Sprintf("%v", exists)},
	}, nil
}

func (p *Plugin) handleDeleteNamespace(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	if err := c.DropNamespace(ctx, namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Namespace %v deleted", namespace),
	}, nil
}

func (p *Plugin) handleUpdateNamespaceProps(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var removals []string
	if removalsJSON, ok := req.Payload["removals"]; ok {
		if err := json.Unmarshal([]byte(removalsJSON), &removals); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid removals format: %v", err),
			}, nil
		}
	}

	var updates map[string]string
	if updatesJSON, ok := req.Payload["updates"]; ok {
		if err := json.Unmarshal([]byte(updatesJSON), &updates); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid updates format: %v", err),
			}, nil
		}
	}

	resp, err := c.UpdateNamespaceProperties(ctx, namespace, removals, updates)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: "Namespace properties updated",
		Result:  map[string]string{"data": string(data)},
	}, nil
}

// ─── TABLE HANDLERS ─────────────────────────────────────────────────────────

func (p *Plugin) handleListTables(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.ListTables(ctx, namespace, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d tables", len(resp.Identifiers)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var createReq types.CreateTableRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreateTable(ctx, namespace, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table %s created", createReq.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	snapshots, _ := req.Payload["snapshots"]

	resp, err := c.LoadTable(ctx, namespace, tableName, snapshots)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded table %s", tableName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleTableExists(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	exists, err := c.TableExists(ctx, namespace, tableName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table exists: %v", exists),
		Result:  map[string]string{"exists": fmt.Sprintf("%v", exists)},
	}, nil
}

func (p *Plugin) handleUpdateTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	var commitReq types.CommitTableRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &commitReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdateTable(ctx, namespace, tableName, commitReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table %s updated", tableName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeleteTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	purgeRequested := req.Payload["purgeRequested"] == "true"

	if err := c.DropTable(ctx, namespace, tableName, purgeRequested); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table %s deleted", tableName),
	}, nil
}

func (p *Plugin) handleRenameTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	sourceJSON, ok := req.Payload["source"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: source",
		}, nil
	}

	var source types.TableIdentifier
	if err := json.Unmarshal([]byte(sourceJSON), &source); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid source format: %v", err),
		}, nil
	}

	destJSON, ok := req.Payload["destination"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: destination",
		}, nil
	}

	var dest types.TableIdentifier
	if err := json.Unmarshal([]byte(destJSON), &dest); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid destination format: %v", err),
		}, nil
	}

	if err := c.RenameTable(ctx, source, dest); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table renamed from %s to %s", source.Name, dest.Name),
	}, nil
}

func (p *Plugin) handleRegisterTable(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var registerReq types.RegisterTableRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &registerReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.RegisterTable(ctx, namespace, registerReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Table %s registered", registerReq.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadTableCreds(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	resp, err := c.LoadCredentials(ctx, namespace, tableName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded credentials for table %s", tableName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleReportMetrics(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	var metricsReq types.ReportMetricsRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &metricsReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.ReportMetrics(ctx, namespace, tableName, metricsReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: "Metrics reported successfully",
	}, nil
}

func (p *Plugin) handleCommitTransaction(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	var txReq types.CommitTransactionRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &txReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.CommitTransaction(ctx, txReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Transaction committed with %d table changes", len(txReq.TableChanges)),
	}, nil
}

// ─── VIEW HANDLERS ───────────────────────────────────────────────────────────

func (p *Plugin) handleListViews(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.ListViews(ctx, namespace, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d views", len(resp.Identifiers)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateView(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var createReq types.CreateViewRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreateView(ctx, namespace, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("View %s created", createReq.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadView(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	viewName, ok := req.Payload["view"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: view",
		}, nil
	}

	resp, err := c.LoadView(ctx, namespace, viewName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded view %s", viewName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleViewExists(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	viewName, ok := req.Payload["view"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: view",
		}, nil
	}

	exists, err := c.ViewExists(ctx, namespace, viewName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("View exists: %v", exists),
		Result:  map[string]string{"exists": fmt.Sprintf("%v", exists)},
	}, nil
}

func (p *Plugin) handleReplaceView(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	viewName, ok := req.Payload["view"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: view",
		}, nil
	}

	var commitReq types.CommitViewRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &commitReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.ReplaceView(ctx, namespace, viewName, commitReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("View %s replaced", viewName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeleteView(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	viewName, ok := req.Payload["view"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: view",
		}, nil
	}

	if err := c.DropView(ctx, namespace, viewName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("View %s deleted", viewName),
	}, nil
}

func (p *Plugin) handleRenameView(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	sourceJSON, ok := req.Payload["source"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: source",
		}, nil
	}

	var source types.TableIdentifier
	if err := json.Unmarshal([]byte(sourceJSON), &source); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid source format: %v", err),
		}, nil
	}

	destJSON, ok := req.Payload["destination"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: destination",
		}, nil
	}

	var dest types.TableIdentifier
	if err := json.Unmarshal([]byte(destJSON), &dest); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid destination format: %v", err),
		}, nil
	}

	if err := c.RenameView(ctx, source, dest); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("View renamed from %s to %s", source.Name, dest.Name),
	}, nil
}

// ─── CATALOG HANDLERS ────────────────────────────────────────────────────────

func (p *Plugin) handleListCatalogs(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	resp, err := c.ListCatalogs(ctx)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d catalogs", len(resp.Catalogs)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateCatalog(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	var createReq types.CreateCatalogRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreateCatalog(ctx, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog %s created", resp.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadCatalog(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	resp, err := c.GetCatalog(ctx, catalogName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded catalog %s", catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleUpdateCatalog(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	var updateReq types.UpdateCatalogRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &updateReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdateCatalog(ctx, catalogName, updateReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog %s updated", catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeleteCatalog(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	if err := c.DeleteCatalog(ctx, catalogName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog %s deleted", catalogName),
	}, nil
}

// ─── PRINCIPAL HANDLERS ─────────────────────────────────────────────────────

func (p *Plugin) handleListPrincipals(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	resp, err := c.ListPrincipals(ctx)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d principals", len(resp.Principals)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreatePrincipal(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	var createReq types.CreatePrincipalRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreatePrincipal(ctx, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal %s created", resp.Principal.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadPrincipal(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	resp, err := c.GetPrincipal(ctx, principalName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded principal %s", principalName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleUpdatePrincipal(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	var updateReq types.UpdatePrincipalRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &updateReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdatePrincipal(ctx, principalName, updateReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal %s updated", principalName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeletePrincipal(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	if err := c.DeletePrincipal(ctx, principalName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal %s deleted", principalName),
	}, nil
}

func (p *Plugin) handleRotateCreds(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	resp, err := c.RotateCredentials(ctx, principalName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Credentials rotated for principal %s", principalName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleResetCreds(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	var resetReq types.ResetPrincipalRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &resetReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	}

	resp, err := c.ResetCredentials(ctx, principalName, resetReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Credentials reset for principal %s", principalName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

// ─── PRINCIPAL ROLE HANDLERS ───────────────────────────────────────────────

func (p *Plugin) handleListPrincipalRoles(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	resp, err := c.ListPrincipalRoles(ctx)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d principal roles", len(resp.Roles)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreatePrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	var createReq types.CreatePrincipalRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreatePrincipalRole(ctx, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal role %s created", resp.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadPrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	resp, err := c.GetPrincipalRole(ctx, roleName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded principal role %s", roleName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleUpdatePrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	var updateReq types.UpdatePrincipalRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &updateReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdatePrincipalRole(ctx, roleName, updateReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal role %s updated", roleName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeletePrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	if err := c.DeletePrincipalRole(ctx, roleName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Principal role %s deleted", roleName),
	}, nil
}

func (p *Plugin) handleAssignPrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	var grantReq types.GrantPrincipalRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &grantReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.AssignPrincipalRole(ctx, principalName, grantReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Role %s assigned to principal %s", grantReq.PrincipalRole.Name, principalName),
	}, nil
}

func (p *Plugin) handleRevokePrincipalRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	if err := c.RevokePrincipalRole(ctx, principalName, roleName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Role %s revoked from principal %s", roleName, principalName),
	}, nil
}

func (p *Plugin) handleListPrincipalRolesAssigned(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalName, ok := req.Payload["principal"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principal",
		}, nil
	}

	resp, err := c.ListPrincipalRolesAssigned(ctx, principalName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d roles assigned to principal %s", len(resp.Roles), principalName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleListPrincipalsByRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	resp, err := c.ListPrincipalsByRole(ctx, roleName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d principals with role %s", len(resp.Principals), roleName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

// ─── CATALOG ROLE HANDLERS ──────────────────────────────────────────────────

func (p *Plugin) handleListCatalogRoles(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	resp, err := c.ListCatalogRoles(ctx, catalogName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d catalog roles in %s", len(resp.Roles), catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	var createReq types.CreateCatalogRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreateCatalogRole(ctx, catalogName, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog role %s created in %s", resp.Name, catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	resp, err := c.GetCatalogRole(ctx, catalogName, roleName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded catalog role %s from %s", roleName, catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleUpdateCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	var updateReq types.UpdateCatalogRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &updateReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdateCatalogRole(ctx, catalogName, roleName, updateReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog role %s updated in %s", roleName, catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeleteCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	roleName, ok := req.Payload["role"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: role",
		}, nil
	}

	if err := c.DeleteCatalogRole(ctx, catalogName, roleName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog role %s deleted from %s", roleName, catalogName),
	}, nil
}

func (p *Plugin) handleAssignCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalRoleName, ok := req.Payload["principalRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principalRole",
		}, nil
	}

	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	var grantReq types.GrantCatalogRoleRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &grantReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.AssignCatalogRoleToPrincipalRole(ctx, principalRoleName, catalogName, grantReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog role %s assigned to principal role %s in catalog %s", grantReq.CatalogRole.Name, principalRoleName, catalogName),
	}, nil
}

func (p *Plugin) handleRevokeCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalRoleName, ok := req.Payload["principalRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principalRole",
		}, nil
	}

	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	catalogRoleName, ok := req.Payload["catalogRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalogRole",
		}, nil
	}

	if err := c.RevokeCatalogRoleFromPrincipalRole(ctx, principalRoleName, catalogName, catalogRoleName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Catalog role %s revoked from principal role %s in catalog %s", catalogRoleName, principalRoleName, catalogName),
	}, nil
}

func (p *Plugin) handleListCatalogRolesByRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	principalRoleName, ok := req.Payload["principalRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: principalRole",
		}, nil
	}

	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	resp, err := c.ListCatalogRolesForPrincipalRole(ctx, principalRoleName, catalogName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d catalog roles for principal role %s in catalog %s", len(resp.Roles), principalRoleName, catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleListPrincipalRolesByCatalogRole(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	catalogRoleName, ok := req.Payload["catalogRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalogRole",
		}, nil
	}

	resp, err := c.ListPrincipalRolesForCatalogRole(ctx, catalogName, catalogRoleName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d principal roles for catalog role %s in catalog %s", len(resp.Roles), catalogRoleName, catalogName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

// ─── GRANT HANDLERS ──────────────────────────────────────────────────────────

func (p *Plugin) handleListGrants(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	catalogRoleName, ok := req.Payload["catalogRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalogRole",
		}, nil
	}

	resp, err := c.ListGrants(ctx, catalogName, catalogRoleName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d grants for catalog role %s", len(resp.Grants), catalogRoleName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleAddGrant(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	catalogRoleName, ok := req.Payload["catalogRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalogRole",
		}, nil
	}

	var grantReq types.AddGrantRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &grantReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.AddGrant(ctx, catalogName, catalogRoleName, grantReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Grant added to catalog role %s", catalogRoleName),
	}, nil
}

func (p *Plugin) handleRevokeGrant(ctx context.Context, c *client.ManagementClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	catalogName, ok := req.Payload["catalog"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalog",
		}, nil
	}

	catalogRoleName, ok := req.Payload["catalogRole"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: catalogRole",
		}, nil
	}

	var revokeReq types.RevokeGrantRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &revokeReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	cascade := req.Payload["cascade"] == "true"

	if err := c.RevokeGrant(ctx, catalogName, catalogRoleName, revokeReq, cascade); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Grant revoked from catalog role %s", catalogRoleName),
	}, nil
}

// ─── POLICY HANDLERS ─────────────────────────────────────────────────────────

func (p *Plugin) handleListPolicies(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyType, _ := req.Payload["policyType"]
	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.ListPolicies(ctx, namespace, policyType, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d policies", len(resp.Identifiers)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreatePolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var createReq types.CreatePolicyRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreatePolicy(ctx, namespace, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Policy %s created", resp.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadPolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyName, ok := req.Payload["policy"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: policy",
		}, nil
	}

	resp, err := c.GetPolicy(ctx, namespace, policyName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded policy %s", policyName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleUpdatePolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyName, ok := req.Payload["policy"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: policy",
		}, nil
	}

	var updateReq types.UpdatePolicyRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &updateReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.UpdatePolicy(ctx, namespace, policyName, updateReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Policy %s updated", policyName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeletePolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyName, ok := req.Payload["policy"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: policy",
		}, nil
	}

	detachAll := req.Payload["detachAll"] == "true"

	if err := c.DeletePolicy(ctx, namespace, policyName, detachAll); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Policy %s deleted", policyName),
	}, nil
}

func (p *Plugin) handleAttachPolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyName, ok := req.Payload["policy"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: policy",
		}, nil
	}

	var attachReq types.AttachPolicyRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &attachReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.AttachPolicy(ctx, namespace, policyName, attachReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Policy %s attached", policyName),
	}, nil
}

func (p *Plugin) handleDetachPolicy(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	policyName, ok := req.Payload["policy"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: policy",
		}, nil
	}

	var detachReq types.DetachPolicyRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &detachReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.DetachPolicy(ctx, namespace, policyName, detachReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Policy %s detached", policyName),
	}, nil
}

func (p *Plugin) handleGetApplicablePolicies(ctx context.Context, c *client.PolicyClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	var namespace []string
	if nsJSON, ok := req.Payload["namespace"]; ok {
		if err := json.Unmarshal([]byte(nsJSON), &namespace); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid namespace format: %v", err),
			}, nil
		}
	}

	targetName, _ := req.Payload["targetName"]
	policyType, _ := req.Payload["policyType"]
	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.GetApplicablePolicies(ctx, namespace, targetName, policyType, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d applicable policies", len(resp.ApplicablePolicies)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

// ─── GENERIC TABLE HANDLERS ──────────────────────────────────────────────────

func (p *Plugin) handleListGenericTables(ctx context.Context, c *client.GenericClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	pageToken, _ := req.Payload["pageToken"]
	pageSize, _ := req.Payload["pageSize"]

	size := 50
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &size)
	}

	resp, err := c.ListGenericTables(ctx, namespace, pageToken, size)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Found %d generic tables", len(resp.Identifiers)),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleCreateGenericTable(ctx context.Context, c *client.GenericClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	var createReq types.CreateGenericTableRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &createReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	resp, err := c.CreateGenericTable(ctx, namespace, createReq)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Generic table %s created", resp.Name),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleLoadGenericTable(ctx context.Context, c *client.GenericClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	resp, err := c.LoadGenericTable(ctx, namespace, tableName)
	if err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(resp)
	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Loaded generic table %s", tableName),
		Result:  map[string]string{"data": string(data)},
	}, nil
}

func (p *Plugin) handleDeleteGenericTable(ctx context.Context, c *client.GenericClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	if err := c.DropGenericTable(ctx, namespace, tableName); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Generic table %s deleted", tableName),
	}, nil
}

// ─── NOTIFICATION HANDLERS ──────────────────────────────────────────────────

func (p *Plugin) handleSendNotification(ctx context.Context, c *client.IcebergClient, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	namespaceJSON, ok := req.Payload["namespace"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: namespace",
		}, nil
	}

	var namespace []string
	if err := json.Unmarshal([]byte(namespaceJSON), &namespace); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("invalid namespace format: %v", err),
		}, nil
	}

	tableName, ok := req.Payload["table"]
	if !ok {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: table",
		}, nil
	}

	var notifReq types.NotificationRequest
	if reqJSON, ok := req.Payload["request"]; ok {
		if err := json.Unmarshal([]byte(reqJSON), &notifReq); err != nil {
			return &pluginv1.ExecuteResponse{
				State:   pluginv1.JobState_JOB_STATE_FAILED,
				Message: fmt.Sprintf("invalid request format: %v", err),
			}, nil
		}
	} else {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: "missing required payload field: request",
		}, nil
	}

	if err := c.SendNotification(ctx, namespace, tableName, notifReq); err != nil {
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &pluginv1.ExecuteResponse{
		State:   pluginv1.JobState_JOB_STATE_SUCCEEDED,
		Message: fmt.Sprintf("Notification sent to table %s", tableName),
	}, nil
}