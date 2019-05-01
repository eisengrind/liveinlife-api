package rbac

// RoleID id of a role
type RoleID string

// RoleRules from a role
type RoleRules []Rule

// Contains checks whether a rule is available in the role rules
func (r RoleRules) Contains(rule Rule) bool {
	for _, v := range r {
		if rule == v {
			return true
		}
	}

	return false
}
