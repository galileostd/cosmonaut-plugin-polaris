package types

// ─── Iceberg REST API Types ──────────────────────────────────────────────────

type CatalogConfig struct {
    Overrides map[string]string `json:"overrides"`
    Defaults  map[string]string `json:"defaults"`
    Endpoints []string          `json:"endpoints,omitempty"`
}

type Namespace struct {
    Parts []string `json:"namespace"`
}

type ListNamespacesResponse struct {
    NextPageToken string      `json:"next-page-token,omitempty"`
    Namespaces    [][]string  `json:"namespaces"`
}

type CreateNamespaceRequest struct {
    Namespace  []string          `json:"namespace"`
    Properties map[string]string `json:"properties,omitempty"`
}

type CreateNamespaceResponse struct {
    Namespace  []string          `json:"namespace"`
    Properties map[string]string `json:"properties,omitempty"`
}

type GetNamespaceResponse struct {
    Namespace  []string          `json:"namespace"`
    Properties map[string]string `json:"properties,omitempty"`
}

type UpdateNamespacePropertiesRequest struct {
    Removals []string          `json:"removals,omitempty"`
    Updates  map[string]string `json:"updates,omitempty"`
}

type UpdateNamespacePropertiesResponse struct {
    Updated []string `json:"updated"`
    Removed []string `json:"removed"`
    Missing []string `json:"missing,omitempty"`
}

type TableIdentifier struct {
    Namespace []string `json:"namespace"`
    Name      string   `json:"name"`
}

type ListTablesResponse struct {
    NextPageToken string             `json:"next-page-token,omitempty"`
    Identifiers   []TableIdentifier  `json:"identifiers"`
}

type CreateTableRequest struct {
    Name          string            `json:"name"`
    Location      string            `json:"location,omitempty"`
    Schema        interface{}       `json:"schema"`
    PartitionSpec interface{}       `json:"partition-spec,omitempty"`
    WriteOrder    interface{}       `json:"write-order,omitempty"`
    StageCreate   bool              `json:"stage-create,omitempty"`
    Properties    map[string]string `json:"properties,omitempty"`
}

type LoadTableResult struct {
    MetadataLocation string                 `json:"metadata-location,omitempty"`
    Metadata         map[string]interface{} `json:"metadata"`
    Config           map[string]string      `json:"config,omitempty"`
}

type CommitTableRequest struct {
    Identifier   TableIdentifier   `json:"identifier,omitempty"`
    Requirements []interface{}     `json:"requirements"`
    Updates      []interface{}     `json:"updates"`
}

type CommitTableResponse struct {
    MetadataLocation string                 `json:"metadata-location"`
    Metadata         map[string]interface{} `json:"metadata"`
}

type RenameTableRequest struct {
    Source      TableIdentifier `json:"source"`
    Destination TableIdentifier `json:"destination"`
}

type RegisterTableRequest struct {
    Name             string `json:"name"`
    MetadataLocation string `json:"metadata-location"`
    Overwrite        bool   `json:"overwrite,omitempty"`
}

type LoadCredentialsResponse struct {
    StorageCredentials []StorageCredential `json:"storage-credentials"`
}

type StorageCredential struct {
    Prefix string            `json:"prefix"`
    Config map[string]string `json:"config"`
}

type ReportMetricsRequest struct {
    ReportType string                 `json:"report-type"`
    TableName  string                 `json:"table-name"`
    SnapshotID int64                  `json:"snapshot-id"`
    Filter     interface{}            `json:"filter"`
    Metrics    map[string]interface{} `json:"metrics"`
}

type CommitTransactionRequest struct {
    TableChanges []CommitTableRequest `json:"table-changes"`
}

// ─── View Types ────────────────────────────────────────────────────────────────

type ViewRepresentation struct {
    Type    string `json:"type"`
    SQL     string `json:"sql"`
    Dialect string `json:"dialect"`
}

type ViewVersion struct {
    VersionID        int64                  `json:"version-id"`
    TimestampMs      int64                  `json:"timestamp-ms"`
    SchemaID         int32                  `json:"schema-id"`
    Summary          map[string]string      `json:"summary"`
    Representations []ViewRepresentation    `json:"representations"`
    DefaultNamespace []string               `json:"default-namespace"`
    DefaultCatalog   string                 `json:"default-catalog,omitempty"`
}

type ViewMetadata struct {
    ViewUUID        string            `json:"view-uuid"`
    FormatVersion   int               `json:"format-version"`
    Location        string            `json:"location"`
    CurrentVersionID int64            `json:"current-version-id"`
    Versions        []ViewVersion     `json:"versions"`
    VersionLog      []ViewHistoryEntry `json:"version-log"`
    Schemas         []interface{}     `json:"schemas"`
    Properties      map[string]string `json:"properties,omitempty"`
}

type ViewHistoryEntry struct {
    VersionID  int64 `json:"version-id"`
    TimestampMs int64 `json:"timestamp-ms"`
}

type CreateViewRequest struct {
    Name       string            `json:"name"`
    Location   string            `json:"location,omitempty"`
    Schema     interface{}       `json:"schema"`
    ViewVersion ViewVersion      `json:"view-version"`
    Properties map[string]string `json:"properties,omitempty"`
}

type LoadViewResult struct {
    MetadataLocation string            `json:"metadata-location"`
    Metadata         ViewMetadata      `json:"metadata"`
    Config           map[string]string `json:"config,omitempty"`
}

type CommitViewRequest struct {
    Identifier   TableIdentifier   `json:"identifier,omitempty"`
    Requirements []interface{}     `json:"requirements,omitempty"`
    Updates      []interface{}     `json:"updates"`
}

// ─── Notification Types ──────────────────────────────────────────────────────

type NotificationType string

const (
    NotificationCreate  NotificationType = "CREATE"
    NotificationUpdate  NotificationType = "UPDATE"
    NotificationDrop    NotificationType = "DROP"
    NotificationValidate NotificationType = "VALIDATE"
)

type TableUpdateNotification struct {
    TableName        string                 `json:"table-name"`
    Timestamp        int64                  `json:"timestamp"`
    TableUUID        string                 `json:"table-uuid"`
    MetadataLocation string                 `json:"metadata-location"`
    Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type NotificationRequest struct {
    NotificationType NotificationType        `json:"notification-type"`
    Payload          TableUpdateNotification `json:"payload"`
}

// ─── Polaris Management API Types ────────────────────────────────────────────

type Catalog struct {
    Type              string              `json:"type"` // INTERNAL, EXTERNAL
    Name              string              `json:"name"`
    Properties        map[string]string   `json:"properties"`
    CreateTimestamp   int64               `json:"createTimestamp,omitempty"`
    LastUpdateTimestamp int64             `json:"lastUpdateTimestamp,omitempty"`
    EntityVersion     int                 `json:"entityVersion,omitempty"`
    StorageConfigInfo StorageConfigInfo   `json:"storageConfigInfo"`
}

type StorageConfigInfo struct {
    StorageType      string   `json:"storageType"` // S3, GCS, AZURE, FILE
    AllowedLocations []string `json:"allowedLocations,omitempty"`
    StorageName      string   `json:"storageName,omitempty"`
    // S3-specific
    RoleArn          string   `json:"roleArn,omitempty"`
    ExternalID       string   `json:"externalId,omitempty"`
    UserArn          string   `json:"userArn,omitempty"`
    Region           string   `json:"region,omitempty"`
    Endpoint         string   `json:"endpoint,omitempty"`
    PathStyleAccess  bool     `json:"pathStyleAccess,omitempty"`
}

type CreateCatalogRequest struct {
    Catalog Catalog `json:"catalog"`
}

type UpdateCatalogRequest struct {
    CurrentEntityVersion int                 `json:"currentEntityVersion"`
    Properties           map[string]string   `json:"properties,omitempty"`
    StorageConfigInfo    *StorageConfigInfo  `json:"storageConfigInfo,omitempty"`
}

type CatalogsResponse struct {
    Catalogs []Catalog `json:"catalogs"`
}

type Principal struct {
    Name              string            `json:"name"`
    ClientID          string            `json:"clientId,omitempty"`
    Properties        map[string]string `json:"properties,omitempty"`
    CreateTimestamp   int64             `json:"createTimestamp,omitempty"`
    LastUpdateTimestamp int64           `json:"lastUpdateTimestamp,omitempty"`
    EntityVersion     int               `json:"entityVersion,omitempty"`
}

type CreatePrincipalRequest struct {
    Principal                Principal `json:"principal"`
    CredentialRotationRequired bool    `json:"credentialRotationRequired,omitempty"`
}

type PrincipalWithCredentials struct {
    Principal   Principal `json:"principal"`
    Credentials struct {
        ClientID     string `json:"clientId"`
        ClientSecret string `json:"clientSecret"`
    } `json:"credentials"`
}

type UpdatePrincipalRequest struct {
    CurrentEntityVersion int               `json:"currentEntityVersion"`
    Properties           map[string]string `json:"properties"`
}

type ResetPrincipalRequest struct {
    ClientID     string `json:"clientId,omitempty"`
    ClientSecret string `json:"clientSecret,omitempty"`
}

type PrincipalsResponse struct {
    Principals []Principal `json:"principals"`
}

type PrincipalRole struct {
    Name              string            `json:"name"`
    Federated         bool              `json:"federated,omitempty"`
    Properties        map[string]string `json:"properties,omitempty"`
    CreateTimestamp   int64             `json:"createTimestamp,omitempty"`
    LastUpdateTimestamp int64           `json:"lastUpdateTimestamp,omitempty"`
    EntityVersion     int               `json:"entityVersion,omitempty"`
}

type CreatePrincipalRoleRequest struct {
    PrincipalRole PrincipalRole `json:"principalRole"`
}

type UpdatePrincipalRoleRequest struct {
    CurrentEntityVersion int               `json:"currentEntityVersion"`
    Properties           map[string]string `json:"properties"`
}

type PrincipalRolesResponse struct {
    Roles []PrincipalRole `json:"roles"`
}

type GrantPrincipalRoleRequest struct {
    PrincipalRole PrincipalRole `json:"principalRole"`
}

type GrantCatalogRoleRequest struct {
    CatalogRole CatalogRole `json:"catalogRole"`
}

type CatalogRole struct {
    Name              string            `json:"name"`
    Properties        map[string]string `json:"properties,omitempty"`
    CreateTimestamp   int64             `json:"createTimestamp,omitempty"`
    LastUpdateTimestamp int64           `json:"lastUpdateTimestamp,omitempty"`
    EntityVersion     int               `json:"entityVersion,omitempty"`
}

type CreateCatalogRoleRequest struct {
    CatalogRole CatalogRole `json:"catalogRole"`
}

type UpdateCatalogRoleRequest struct {
    CurrentEntityVersion int               `json:"currentEntityVersion"`
    Properties           map[string]string `json:"properties"`
}

type CatalogRolesResponse struct {
    Roles []CatalogRole `json:"roles"`
}

// ─── Grant Types ──────────────────────────────────────────────────────────────

type GrantResourceType string

const (
    GrantTypeCatalog   GrantResourceType = "catalog"
    GrantTypeNamespace GrantResourceType = "namespace"
    GrantTypeTable     GrantResourceType = "table"
    GrantTypeView      GrantResourceType = "view"
    GrantTypePolicy    GrantResourceType = "policy"
)

type GrantResource struct {
    Type GrantResourceType `json:"type"`
    // Catalog
    Privilege string `json:"privilege,omitempty"`
    // Namespace
    Namespace []string `json:"namespace,omitempty"`
    // Table
    TableName string `json:"tableName,omitempty"`
    // View
    ViewName string `json:"viewName,omitempty"`
    // Policy
    PolicyName string `json:"policyName,omitempty"`
}

type AddGrantRequest struct {
    Grant GrantResource `json:"grant"`
}

type RevokeGrantRequest struct {
    Grant GrantResource `json:"grant"`
}

type GrantResourcesResponse struct {
    Grants []GrantResource `json:"grants"`
}

// ─── Policy API Types ─────────────────────────────────────────────────────────

type PolicyType string

type Policy struct {
    PolicyType  string `json:"policy-type"`
    Inheritable bool   `json:"inheritable"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Content     string `json:"content,omitempty"`
    Version     int    `json:"version"`
}

type CreatePolicyRequest struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description,omitempty"`
    Content     string `json:"content,omitempty"`
}

type UpdatePolicyRequest struct {
    Description          string `json:"description,omitempty"`
    Content              string `json:"content,omitempty"`
    CurrentPolicyVersion int    `json:"current-policy-version"`
}

type PolicyIdentifier struct {
    Namespace []string `json:"namespace"`
    Name      string   `json:"name"`
}

type ListPoliciesResponse struct {
    NextPageToken string             `json:"next-page-token,omitempty"`
    Identifiers   []PolicyIdentifier `json:"identifiers"`
}

type PolicyAttachmentTarget struct {
    Type string   `json:"type"` // catalog, namespace, table-like
    Path []string `json:"path,omitempty"`
}

type AttachPolicyRequest struct {
    Target     PolicyAttachmentTarget `json:"target"`
    Parameters map[string]string      `json:"parameters,omitempty"`
}

type DetachPolicyRequest struct {
    Target PolicyAttachmentTarget `json:"target"`
}

type ApplicablePolicy struct {
    Policy
    Inherited bool     `json:"inherited"`
    Namespace []string `json:"namespace"`
}

type GetApplicablePoliciesResponse struct {
    NextPageToken      string              `json:"next-page-token,omitempty"`
    ApplicablePolicies []ApplicablePolicy  `json:"applicable-policies"`
}

// ─── Generic Table Types ──────────────────────────────────────────────────────

type GenericTable struct {
    Name         string            `json:"name"`
    Format       string            `json:"format"`
    BaseLocation string            `json:"base-location,omitempty"`
    Doc          string            `json:"doc,omitempty"`
    Properties   map[string]string `json:"properties,omitempty"`
}

type CreateGenericTableRequest struct {
    Name         string            `json:"name"`
    Format       string            `json:"format"`
    BaseLocation string            `json:"base-location,omitempty"`
    Doc          string            `json:"doc,omitempty"`
    Properties   map[string]string `json:"properties,omitempty"`
}

type LoadGenericTableResponse struct {
    Table GenericTable `json:"table"`
}

type ListGenericTablesResponse struct {
    NextPageToken string             `json:"next-page-token,omitempty"`
    Identifiers   []TableIdentifier  `json:"identifiers"`
}