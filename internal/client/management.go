package client

import (
	"context"
	"net/url"

	"github.com/galileostd/cosmonaut-plugin-polaris/internal/types"
)

// ManagementClient handles Polaris Management API operations.
type ManagementClient struct {
	*PolarisClient
}

// NewManagementClient creates a new Management client.
func NewManagementClient(cfg PolarisClientConfig) *ManagementClient {
	return &ManagementClient{PolarisClient: NewPolarisClient(cfg)}
}

// ─── Catalogs ──────────────────────────────────────────────────────────────────

// ListCatalogs lists all catalogs.
func (c *ManagementClient) ListCatalogs(ctx context.Context) (*types.CatalogsResponse, error) {
	url := c.ManagementURL("catalogs")
	var resp types.CatalogsResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCatalog creates a new catalog.
func (c *ManagementClient) CreateCatalog(ctx context.Context, req types.CreateCatalogRequest) (*types.Catalog, error) {
	url := c.ManagementURL("catalogs")
	var resp types.Catalog
	if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCatalog loads a catalog by name.
func (c *ManagementClient) GetCatalog(ctx context.Context, catalogName string) (*types.Catalog, error) {
	url := c.ManagementURL("catalogs", catalogName)
	var resp types.Catalog
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCatalog updates an existing catalog.
func (c *ManagementClient) UpdateCatalog(ctx context.Context, catalogName string, req types.UpdateCatalogRequest) (*types.Catalog, error) {
	url := c.ManagementURL("catalogs", catalogName)
	var resp types.Catalog
	if err := c.doRequest(ctx, "PUT", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCatalog deletes a catalog (must be empty).
func (c *ManagementClient) DeleteCatalog(ctx context.Context, catalogName string) error {
	url := c.ManagementURL("catalogs", catalogName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ─── Principals ────────────────────────────────────────────────────────────────

// ListPrincipals lists all principals.
func (c *ManagementClient) ListPrincipals(ctx context.Context) (*types.PrincipalsResponse, error) {
	url := c.ManagementURL("principals")
	var resp types.PrincipalsResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreatePrincipal creates a new principal.
func (c *ManagementClient) CreatePrincipal(ctx context.Context, req types.CreatePrincipalRequest) (*types.PrincipalWithCredentials, error) {
	url := c.ManagementURL("principals")
	var resp types.PrincipalWithCredentials
	if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPrincipal loads a principal by name.
func (c *ManagementClient) GetPrincipal(ctx context.Context, principalName string) (*types.Principal, error) {
	url := c.ManagementURL("principals", principalName)
	var resp types.Principal
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatePrincipal updates an existing principal.
func (c *ManagementClient) UpdatePrincipal(ctx context.Context, principalName string, req types.UpdatePrincipalRequest) (*types.Principal, error) {
	url := c.ManagementURL("principals", principalName)
	var resp types.Principal
	if err := c.doRequest(ctx, "PUT", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeletePrincipal deletes a principal.
func (c *ManagementClient) DeletePrincipal(ctx context.Context, principalName string) error {
	url := c.ManagementURL("principals", principalName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// RotateCredentials rotates a principal's credentials.
func (c *ManagementClient) RotateCredentials(ctx context.Context, principalName string) (*types.PrincipalWithCredentials, error) {
	url := c.ManagementURL("principals", principalName, "rotate")
	var resp types.PrincipalWithCredentials
	if err := c.doRequest(ctx, "POST", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetCredentials resets a principal's credentials.
func (c *ManagementClient) ResetCredentials(ctx context.Context, principalName string, req types.ResetPrincipalRequest) (*types.PrincipalWithCredentials, error) {
	url := c.ManagementURL("principals", principalName, "reset")
	var resp types.PrincipalWithCredentials
	if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ─── Principal Roles ──────────────────────────────────────────────────────────

// ListPrincipalRoles lists all principal roles.
func (c *ManagementClient) ListPrincipalRoles(ctx context.Context) (*types.PrincipalRolesResponse, error) {
	url := c.ManagementURL("principal-roles")
	var resp types.PrincipalRolesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreatePrincipalRole creates a new principal role.
func (c *ManagementClient) CreatePrincipalRole(ctx context.Context, req types.CreatePrincipalRoleRequest) (*types.PrincipalRole, error) {
	url := c.ManagementURL("principal-roles")
	var resp types.PrincipalRole
	if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPrincipalRole loads a principal role by name.
func (c *ManagementClient) GetPrincipalRole(ctx context.Context, roleName string) (*types.PrincipalRole, error) {
	url := c.ManagementURL("principal-roles", roleName)
	var resp types.PrincipalRole
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatePrincipalRole updates an existing principal role.
func (c *ManagementClient) UpdatePrincipalRole(ctx context.Context, roleName string, req types.UpdatePrincipalRoleRequest) (*types.PrincipalRole, error) {
	url := c.ManagementURL("principal-roles", roleName)
	var resp types.PrincipalRole
	if err := c.doRequest(ctx, "PUT", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeletePrincipalRole deletes a principal role.
func (c *ManagementClient) DeletePrincipalRole(ctx context.Context, roleName string) error {
	url := c.ManagementURL("principal-roles", roleName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ─── Principal Role Assignments ──────────────────────────────────────────────

// AssignPrincipalRole assigns a role to a principal.
func (c *ManagementClient) AssignPrincipalRole(ctx context.Context, principalName string, req types.GrantPrincipalRoleRequest) error {
	url := c.ManagementURL("principals", principalName, "principal-roles")
	return c.doRequest(ctx, "PUT", url, req, nil)
}

// RevokePrincipalRole revokes a role from a principal.
func (c *ManagementClient) RevokePrincipalRole(ctx context.Context, principalName, roleName string) error {
	url := c.ManagementURL("principals", principalName, "principal-roles", roleName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ListPrincipalRolesAssigned lists roles assigned to a principal.
func (c *ManagementClient) ListPrincipalRolesAssigned(ctx context.Context, principalName string) (*types.PrincipalRolesResponse, error) {
	url := c.ManagementURL("principals", principalName, "principal-roles")
	var resp types.PrincipalRolesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPrincipalsByRole lists principals assigned to a role.
func (c *ManagementClient) ListPrincipalsByRole(ctx context.Context, roleName string) (*types.PrincipalsResponse, error) {
	url := c.ManagementURL("principal-roles", roleName, "principals")
	var resp types.PrincipalsResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ─── Catalog Roles ────────────────────────────────────────────────────────────

// ListCatalogRoles lists catalog roles in a catalog.
func (c *ManagementClient) ListCatalogRoles(ctx context.Context, catalogName string) (*types.CatalogRolesResponse, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles")
	var resp types.CatalogRolesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCatalogRole creates a new catalog role.
func (c *ManagementClient) CreateCatalogRole(ctx context.Context, catalogName string, req types.CreateCatalogRoleRequest) (*types.CatalogRole, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles")
	var resp types.CatalogRole
	if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCatalogRole loads a catalog role by name.
func (c *ManagementClient) GetCatalogRole(ctx context.Context, catalogName, roleName string) (*types.CatalogRole, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", roleName)
	var resp types.CatalogRole
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCatalogRole updates an existing catalog role.
func (c *ManagementClient) UpdateCatalogRole(ctx context.Context, catalogName, roleName string, req types.UpdateCatalogRoleRequest) (*types.CatalogRole, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", roleName)
	var resp types.CatalogRole
	if err := c.doRequest(ctx, "PUT", url, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCatalogRole deletes a catalog role.
func (c *ManagementClient) DeleteCatalogRole(ctx context.Context, catalogName, roleName string) error {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", roleName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ─── Catalog Role ↔ Principal Role Mapping ──────────────────────────────────

// AssignCatalogRoleToPrincipalRole assigns a catalog role to a principal role.
func (c *ManagementClient) AssignCatalogRoleToPrincipalRole(ctx context.Context, principalRoleName, catalogName string, req types.GrantCatalogRoleRequest) error {
	url := c.ManagementURL("principal-roles", principalRoleName, "catalog-roles", catalogName)
	return c.doRequest(ctx, "PUT", url, req, nil)
}

// RevokeCatalogRoleFromPrincipalRole revokes a catalog role from a principal role.
func (c *ManagementClient) RevokeCatalogRoleFromPrincipalRole(ctx context.Context, principalRoleName, catalogName, catalogRoleName string) error {
	url := c.ManagementURL("principal-roles", principalRoleName, "catalog-roles", catalogName, catalogRoleName)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ListCatalogRolesForPrincipalRole lists catalog roles mapped to a principal role.
func (c *ManagementClient) ListCatalogRolesForPrincipalRole(ctx context.Context, principalRoleName, catalogName string) (*types.CatalogRolesResponse, error) {
	url := c.ManagementURL("principal-roles", principalRoleName, "catalog-roles", catalogName)
	var resp types.CatalogRolesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPrincipalRolesForCatalogRole lists principal roles assigned to a catalog role.
func (c *ManagementClient) ListPrincipalRolesForCatalogRole(ctx context.Context, catalogName, catalogRoleName string) (*types.PrincipalRolesResponse, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", catalogRoleName, "principal-roles")
	var resp types.PrincipalRolesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ─── Grants ───────────────────────────────────────────────────────────────────

// ListGrants lists all grants for a catalog role.
func (c *ManagementClient) ListGrants(ctx context.Context, catalogName, catalogRoleName string) (*types.GrantResourcesResponse, error) {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", catalogRoleName, "grants")
	var resp types.GrantResourcesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddGrant adds a grant to a catalog role.
func (c *ManagementClient) AddGrant(ctx context.Context, catalogName, catalogRoleName string, req types.AddGrantRequest) error {
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", catalogRoleName, "grants")
	return c.doRequest(ctx, "PUT", url, req, nil)
}

// RevokeGrant revokes a grant from a catalog role.
func (c *ManagementClient) RevokeGrant(ctx context.Context, catalogName, catalogRoleName string, req types.RevokeGrantRequest, cascade bool) error {
	params := url.Values{}
	if cascade {
		params.Set("cascade", "true")
	}
	url := c.ManagementURL("catalogs", catalogName, "catalog-roles", catalogRoleName, "grants")
	if len(params) > 0 {
		url = url + "?" + params.Encode()
	}
	return c.doRequest(ctx, "POST", url, req, nil)
}
