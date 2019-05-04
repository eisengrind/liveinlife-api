package rbac

// AccountID of a user in the RBAC-system
type AccountID string

// AccountRoles are the roles attached to a subject
type AccountRoles []RoleID

// Contains checks whether a rule is available in the role rules
func (r AccountRoles) Contains(roleID RoleID) bool {
	for _, v := range r {
		if roleID == v {
			return true
		}
	}

	return false
}
