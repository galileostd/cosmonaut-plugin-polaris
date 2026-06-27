package plugin

import (
	"context"
	"fmt"

	pluginv1 "github.com/galileostd/cosmonaut-sdk/go/plugin/v1"
	sdkserver "github.com/galileostd/cosmonaut-sdk/go/server"

	"github.com/galileostd/cosmonaut-plugin-polaris/internal/client"
)

// Plugin implements the Cosmonaut PluginService for Apache Polaris.
type Plugin struct {
	sdkserver.UnimplementedPlugin
}

// New creates a new Polaris plugin instance.
func New() *Plugin {
	return &Plugin{}
}

// Describe returns static metadata about this plugin.
func (p *Plugin) Describe(_ context.Context, _ *pluginv1.DescribeRequest) (*pluginv1.DescribeResponse, error) {
	return &pluginv1.DescribeResponse{
		PluginName:    "polaris",
		DisplayName:   "Apache Polaris",
		Version:       "v0.1.0",
		Description:   "Enterprise-grade Iceberg catalog with RBAC, policies, and multi-catalog support. Implements Iceberg REST Catalog API + Polaris Management API.",
		PluginType:    pluginv1.PluginType_PLUGIN_TYPE_CATALOG,
		ExecutionType: pluginv1.ExecutionType_EXECUTION_TYPE_OBSERVABILITY,
		Capabilities: []*pluginv1.Capability{
			// Namespaces
			{Type: "list-namespaces", Description: "List all namespaces."},
			{Type: "create-namespace", Description: "Create a new namespace."},
			{Type: "load-namespace", Description: "Load namespace properties."},
			{Type: "namespace-exists", Description: "Check if a namespace exists."},
			{Type: "delete-namespace", Description: "Delete a namespace (must be empty)."},
			{Type: "update-namespace-props", Description: "Update/remove namespace properties."},
			// Tables
			{Type: "list-tables", Description: "List all tables in a namespace."},
			{Type: "create-table", Description: "Create a new table."},
			{Type: "load-table", Description: "Load table metadata."},
			{Type: "table-exists", Description: "Check if a table exists."},
			{Type: "update-table", Description: "Commit updates to a table."},
			{Type: "delete-table", Description: "Delete a table."},
			{Type: "rename-table", Description: "Rename a table."},
			{Type: "register-table", Description: "Register a table from metadata location."},
			{Type: "load-table-creds", Description: "Load vended credentials for a table."},
			{Type: "report-metrics", Description: "Send metrics report."},
			{Type: "commit-transaction", Description: "Atomic commit of multiple tables."},
			// Views
			{Type: "list-views", Description: "List all views in a namespace."},
			{Type: "create-view", Description: "Create a new view."},
			{Type: "load-view", Description: "Load view metadata."},
			{Type: "view-exists", Description: "Check if a view exists."},
			{Type: "replace-view", Description: "Replace a view (commit updates)."},
			{Type: "delete-view", Description: "Delete a view."},
			{Type: "rename-view", Description: "Rename a view."},
			// Catalogs (Management)
			{Type: "list-catalogs", Description: "List all catalogs."},
			{Type: "create-catalog", Description: "Create a new catalog."},
			{Type: "load-catalog", Description: "Load catalog details."},
			{Type: "update-catalog", Description: "Update catalog properties."},
			{Type: "delete-catalog", Description: "Delete a catalog (must be empty)."},
			// Principals
			{Type: "list-principals", Description: "List all principals."},
			{Type: "create-principal", Description: "Create a new principal."},
			{Type: "load-principal", Description: "Load principal details."},
			{Type: "update-principal", Description: "Update principal properties."},
			{Type: "delete-principal", Description: "Delete a principal."},
			{Type: "rotate-creds", Description: "Rotate principal credentials."},
			{Type: "reset-creds", Description: "Reset principal credentials."},
			// Principal Roles
			{Type: "list-principal-roles", Description: "List all principal roles."},
			{Type: "create-principal-role", Description: "Create a new principal role."},
			{Type: "load-principal-role", Description: "Load principal role details."},
			{Type: "update-principal-role", Description: "Update principal role properties."},
			{Type: "delete-principal-role", Description: "Delete a principal role."},
			{Type: "assign-principal-role", Description: "Assign a role to a principal."},
			{Type: "revoke-principal-role", Description: "Revoke a role from a principal."},
			{Type: "list-principal-roles-assigned", Description: "List roles assigned to a principal."},
			{Type: "list-principals-by-role", Description: "List principals assigned to a role."},
			// Catalog Roles
			{Type: "list-catalog-roles", Description: "List catalog roles in a catalog."},
			{Type: "create-catalog-role", Description: "Create a new catalog role."},
			{Type: "load-catalog-role", Description: "Load catalog role details."},
			{Type: "update-catalog-role", Description: "Update catalog role properties."},
			{Type: "delete-catalog-role", Description: "Delete a catalog role."},
			{Type: "assign-catalog-role", Description: "Assign catalog role to principal role."},
			{Type: "revoke-catalog-role", Description: "Revoke catalog role from principal role."},
			{Type: "list-catalog-roles-by-role", Description: "List catalog roles of a principal role."},
			{Type: "list-principal-roles-by-catalog-role", Description: "List principal roles of a catalog role."},
			// Grants
			{Type: "list-grants", Description: "List grants of a catalog role."},
			{Type: "add-grant", Description: "Add a grant to a catalog role."},
			{Type: "revoke-grant", Description: "Revoke a grant from a catalog role."},
			// Policies
			{Type: "list-policies", Description: "List policies in a namespace."},
			{Type: "create-policy", Description: "Create a new policy."},
			{Type: "load-policy", Description: "Load policy details."},
			{Type: "update-policy", Description: "Update policy (version check)."},
			{Type: "delete-policy", Description: "Delete a policy."},
			{Type: "attach-policy", Description: "Attach policy to a target."},
			{Type: "detach-policy", Description: "Detach policy from a target."},
			{Type: "get-applicable-policies", Description: "Get policies applicable to an entity."},
			// Generic Tables
			{Type: "list-generic-tables", Description: "List generic tables in a namespace."},
			{Type: "create-generic-table", Description: "Create a new generic table."},
			{Type: "load-generic-table", Description: "Load generic table details."},
			{Type: "delete-generic-table", Description: "Delete a generic table."},
			// Notifications
			{Type: "send-notification", Description: "Send a notification to a table."},
		},
	}, nil
}

// HealthCheck verifies that Polaris is reachable.
func (p *Plugin) HealthCheck(ctx context.Context, req *pluginv1.HealthCheckRequest) (*pluginv1.HealthCheckResponse, error) {
	if req.Component == nil {
		return unhealthy("component is required"), nil
	}

	icebergClient := client.NewIcebergClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	_, err := icebergClient.GetConfig(ctx)
	if err != nil {
		return &pluginv1.HealthCheckResponse{
			State:   pluginv1.HealthState_HEALTH_STATE_UNHEALTHY,
			Message: fmt.Sprintf("failed to reach Polaris at %s: %v", req.Component.Endpoint, err),
			Details: map[string]string{
				"endpoint": req.Component.Endpoint,
			},
		}, nil
	}

	managementClient := client.NewManagementClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	_, err = managementClient.ListCatalogs(ctx)
	hasManagement := err == nil

	return &pluginv1.HealthCheckResponse{
		State:   pluginv1.HealthState_HEALTH_STATE_HEALTHY,
		Message: "Polaris is reachable and responding",
		Details: map[string]string{
			"endpoint":   req.Component.Endpoint,
			"catalog":    req.Component.Config["catalog"],
			"management": fmt.Sprintf("%v", hasManagement),
		},
	}, nil
}

// Execute runs an action on Polaris.
func (p *Plugin) Execute(ctx context.Context, req *pluginv1.ExecuteRequest) (*pluginv1.ExecuteResponse, error) {
	if req.Component == nil {
		return nil, fmt.Errorf("component is required")
	}

	icebergClient := client.NewIcebergClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	managementClient := client.NewManagementClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	policyClient := client.NewPolicyClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	genericClient := client.NewGenericClient(client.PolarisClientConfig{
		Endpoint: req.Component.Config["endpoint"],
		Token:    req.Component.Config["token"],
		Catalog:  req.Component.Config["catalog"],
		Prefix:   req.Component.Config["prefix"],
	})

	switch req.Action {
	// Namespaces
	case "list-namespaces":
		return p.handleListNamespaces(ctx, icebergClient, req)
	case "create-namespace":
		return p.handleCreateNamespace(ctx, icebergClient, req)
	case "load-namespace":
		return p.handleLoadNamespace(ctx, icebergClient, req)
	case "namespace-exists":
		return p.handleNamespaceExists(ctx, icebergClient, req)
	case "delete-namespace":
		return p.handleDeleteNamespace(ctx, icebergClient, req)
	case "update-namespace-props":
		return p.handleUpdateNamespaceProps(ctx, icebergClient, req)

	// Tables
	case "list-tables":
		return p.handleListTables(ctx, icebergClient, req)
	case "create-table":
		return p.handleCreateTable(ctx, icebergClient, req)
	case "load-table":
		return p.handleLoadTable(ctx, icebergClient, req)
	case "table-exists":
		return p.handleTableExists(ctx, icebergClient, req)
	case "update-table":
		return p.handleUpdateTable(ctx, icebergClient, req)
	case "delete-table":
		return p.handleDeleteTable(ctx, icebergClient, req)
	case "rename-table":
		return p.handleRenameTable(ctx, icebergClient, req)
	case "register-table":
		return p.handleRegisterTable(ctx, icebergClient, req)
	case "load-table-creds":
		return p.handleLoadTableCreds(ctx, icebergClient, req)
	case "report-metrics":
		return p.handleReportMetrics(ctx, icebergClient, req)
	case "commit-transaction":
		return p.handleCommitTransaction(ctx, icebergClient, req)

	// Views
	case "list-views":
		return p.handleListViews(ctx, icebergClient, req)
	case "create-view":
		return p.handleCreateView(ctx, icebergClient, req)
	case "load-view":
		return p.handleLoadView(ctx, icebergClient, req)
	case "view-exists":
		return p.handleViewExists(ctx, icebergClient, req)
	case "replace-view":
		return p.handleReplaceView(ctx, icebergClient, req)
	case "delete-view":
		return p.handleDeleteView(ctx, icebergClient, req)
	case "rename-view":
		return p.handleRenameView(ctx, icebergClient, req)

	// Catalogs
	case "list-catalogs":
		return p.handleListCatalogs(ctx, managementClient, req)
	case "create-catalog":
		return p.handleCreateCatalog(ctx, managementClient, req)
	case "load-catalog":
		return p.handleLoadCatalog(ctx, managementClient, req)
	case "update-catalog":
		return p.handleUpdateCatalog(ctx, managementClient, req)
	case "delete-catalog":
		return p.handleDeleteCatalog(ctx, managementClient, req)

	// Principals
	case "list-principals":
		return p.handleListPrincipals(ctx, managementClient, req)
	case "create-principal":
		return p.handleCreatePrincipal(ctx, managementClient, req)
	case "load-principal":
		return p.handleLoadPrincipal(ctx, managementClient, req)
	case "update-principal":
		return p.handleUpdatePrincipal(ctx, managementClient, req)
	case "delete-principal":
		return p.handleDeletePrincipal(ctx, managementClient, req)
	case "rotate-creds":
		return p.handleRotateCreds(ctx, managementClient, req)
	case "reset-creds":
		return p.handleResetCreds(ctx, managementClient, req)

	// Principal Roles
	case "list-principal-roles":
		return p.handleListPrincipalRoles(ctx, managementClient, req)
	case "create-principal-role":
		return p.handleCreatePrincipalRole(ctx, managementClient, req)
	case "load-principal-role":
		return p.handleLoadPrincipalRole(ctx, managementClient, req)
	case "update-principal-role":
		return p.handleUpdatePrincipalRole(ctx, managementClient, req)
	case "delete-principal-role":
		return p.handleDeletePrincipalRole(ctx, managementClient, req)
	case "assign-principal-role":
		return p.handleAssignPrincipalRole(ctx, managementClient, req)
	case "revoke-principal-role":
		return p.handleRevokePrincipalRole(ctx, managementClient, req)
	case "list-principal-roles-assigned":
		return p.handleListPrincipalRolesAssigned(ctx, managementClient, req)
	case "list-principals-by-role":
		return p.handleListPrincipalsByRole(ctx, managementClient, req)

	// Catalog Roles
	case "list-catalog-roles":
		return p.handleListCatalogRoles(ctx, managementClient, req)
	case "create-catalog-role":
		return p.handleCreateCatalogRole(ctx, managementClient, req)
	case "load-catalog-role":
		return p.handleLoadCatalogRole(ctx, managementClient, req)
	case "update-catalog-role":
		return p.handleUpdateCatalogRole(ctx, managementClient, req)
	case "delete-catalog-role":
		return p.handleDeleteCatalogRole(ctx, managementClient, req)
	case "assign-catalog-role":
		return p.handleAssignCatalogRole(ctx, managementClient, req)
	case "revoke-catalog-role":
		return p.handleRevokeCatalogRole(ctx, managementClient, req)
	case "list-catalog-roles-by-role":
		return p.handleListCatalogRolesByRole(ctx, managementClient, req)
	case "list-principal-roles-by-catalog-role":
		return p.handleListPrincipalRolesByCatalogRole(ctx, managementClient, req)

	// Grants
	case "list-grants":
		return p.handleListGrants(ctx, managementClient, req)
	case "add-grant":
		return p.handleAddGrant(ctx, managementClient, req)
	case "revoke-grant":
		return p.handleRevokeGrant(ctx, managementClient, req)

	// Policies
	case "list-policies":
		return p.handleListPolicies(ctx, policyClient, req)
	case "create-policy":
		return p.handleCreatePolicy(ctx, policyClient, req)
	case "load-policy":
		return p.handleLoadPolicy(ctx, policyClient, req)
	case "update-policy":
		return p.handleUpdatePolicy(ctx, policyClient, req)
	case "delete-policy":
		return p.handleDeletePolicy(ctx, policyClient, req)
	case "attach-policy":
		return p.handleAttachPolicy(ctx, policyClient, req)
	case "detach-policy":
		return p.handleDetachPolicy(ctx, policyClient, req)
	case "get-applicable-policies":
		return p.handleGetApplicablePolicies(ctx, policyClient, req)

	// Generic Tables
	case "list-generic-tables":
		return p.handleListGenericTables(ctx, genericClient, req)
	case "create-generic-table":
		return p.handleCreateGenericTable(ctx, genericClient, req)
	case "load-generic-table":
		return p.handleLoadGenericTable(ctx, genericClient, req)
	case "delete-generic-table":
		return p.handleDeleteGenericTable(ctx, genericClient, req)

	// Notifications
	case "send-notification":
		return p.handleSendNotification(ctx, icebergClient, req)

	default:
		return &pluginv1.ExecuteResponse{
			State:   pluginv1.JobState_JOB_STATE_FAILED,
			Message: fmt.Sprintf("unsupported action: %s", req.Action),
		}, nil
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func unhealthy(msg string) *pluginv1.HealthCheckResponse {
	return &pluginv1.HealthCheckResponse{
		State:   pluginv1.HealthState_HEALTH_STATE_UNHEALTHY,
		Message: msg,
	}
}

// GetLogs returns the Polaris service logs (the Polaris pod). Polaris has no
// per-job pods, so JobId is ignored and we always return the Polaris pod's logs.
func (p *Plugin) GetLogs(ctx context.Context, req *pluginv1.GetLogsRequest) (*pluginv1.GetLogsResponse, error) {
	if req.Component == nil {
		return nil, fmt.Errorf("component is required")
	}

	k8s, err := client.NewK8sClient()
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	namespace := req.Component.Config["service_namespace"]
	selector := req.Component.Config["pod_selector"]

	lines, err := k8s.GetServiceLogs(ctx, namespace, selector, int64(req.TailLines))
	if err != nil {
		return nil, err
	}

	return &pluginv1.GetLogsResponse{Lines: lines}, nil
}