package schema

// ConversationMember - Conversation member info
type ConversationMember struct {

	// Member ID
	ID string `json:"id,omitempty"`

	// Display friendly name
	Name string `json:"name,omitempty"`

	// GivenName
	GivenName string `json:"givenName,omitempty"`

	// Surname
	Surname string `json:"surname,omitempty"`

	// This account's object ID within Azure Active Directory (AAD)
	AadObjectID string `json:"aadObjectId,omitempty"`

	// Email is user email
	Email string `json:"email,omitempty"`

	// UserPrincipalName
	UserPrincipalName string `json:"userPrincipalName,omitempty"`

	// TenantID is an ID fo a tenant
	TenantID string `json:"tenantId,omitempty"`

	// UserRole is a user role
	UserRole RoleTypes `json:"userRole,omitempty"`
}
