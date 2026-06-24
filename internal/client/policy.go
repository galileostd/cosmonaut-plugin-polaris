package client

import (
    "context"
    "net/url"
	"fmt"

    "github.com/galileostd/cosmonaut-plugin-polaris/internal/types"
)

// PolicyClient handles Polaris Policy API operations.
type PolicyClient struct {
    *PolarisClient
}

// NewPolicyClient creates a new Policy client.
func NewPolicyClient(cfg PolarisClientConfig) *PolicyClient {
    return &PolicyClient{PolarisClient: NewPolarisClient(cfg)}
}

// ─── Policies ──────────────────────────────────────────────────────────────────

// ListPolicies lists all policies in a namespace.
func (c *PolicyClient) ListPolicies(ctx context.Context, namespace []string, policyType string, pageToken string, pageSize int) (*types.ListPoliciesResponse, error) {
    params := url.Values{}
    if policyType != "" {
        params.Set("policyType", policyType)
    }
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.ListPoliciesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// CreatePolicy creates a new policy.
func (c *PolicyClient) CreatePolicy(ctx context.Context, namespace []string, req types.CreatePolicyRequest) (*types.Policy, error) {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies")
    var resp types.Policy
    if err := c.doRequest(ctx, "POST", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// GetPolicy loads a policy by name.
func (c *PolicyClient) GetPolicy(ctx context.Context, namespace []string, policyName string) (*types.Policy, error) {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies", policyName)
    var resp types.Policy
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// UpdatePolicy updates an existing policy.
func (c *PolicyClient) UpdatePolicy(ctx context.Context, namespace []string, policyName string, req types.UpdatePolicyRequest) (*types.Policy, error) {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies", policyName)
    var resp types.Policy
    if err := c.doRequest(ctx, "PUT", url, req, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

// DeletePolicy deletes a policy.
func (c *PolicyClient) DeletePolicy(ctx context.Context, namespace []string, policyName string, detachAll bool) error {
    params := url.Values{}
    if detachAll {
        params.Set("detach-all", "true")
    }
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies", policyName)
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }
    return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ─── Policy Attachments ───────────────────────────────────────────────────────

// AttachPolicy attaches a policy to a target.
func (c *PolicyClient) AttachPolicy(ctx context.Context, namespace []string, policyName string, req types.AttachPolicyRequest) error {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies", policyName, "mappings")
    return c.doRequest(ctx, "PUT", url, req, nil)
}

// DetachPolicy detaches a policy from a target.
func (c *PolicyClient) DetachPolicy(ctx context.Context, namespace []string, policyName string, req types.DetachPolicyRequest) error {
    url := c.PolarisURL("namespaces", encodeNamespace(namespace), "policies", policyName, "mappings")
    return c.doRequest(ctx, "POST", url, req, nil)
}

// ─── Applicable Policies ──────────────────────────────────────────────────────

// GetApplicablePolicies gets policies applicable to an entity.
func (c *PolicyClient) GetApplicablePolicies(ctx context.Context, namespace []string, targetName string, policyType string, pageToken string, pageSize int) (*types.GetApplicablePoliciesResponse, error) {
    params := url.Values{}
    if namespace != nil && len(namespace) > 0 {
        params.Set("namespace", encodeNamespace(namespace))
    }
    if targetName != "" {
        params.Set("target-name", targetName)
    }
    if policyType != "" {
        params.Set("policyType", policyType)
    }
    if pageToken != "" {
        params.Set("pageToken", pageToken)
    }
    if pageSize > 0 {
        params.Set("pageSize", fmt.Sprintf("%d", pageSize))
    }

    url := c.PolarisURL("applicable-policies")
    if len(params) > 0 {
        url = url + "?" + params.Encode()
    }

    var resp types.GetApplicablePoliciesResponse
    if err := c.doRequest(ctx, "GET", url, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}