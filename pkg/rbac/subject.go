package rbac

// SubjectID of a subject in the RBAC-system
type SubjectID string

// SubjectRoles are the roles attached to a subject
type SubjectRoles []RoleID

// Contains checks whether a rule is available in the role rules
func (r SubjectRoles) Contains(roleID RoleID) bool {
	for _, v := range r {
		if roleID == v {
			return true
		}
	}

	return false
}
