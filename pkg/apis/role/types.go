package role

import "github.com/51st-state/api/pkg/rbac"

// Identifier of a role object
//go:generate counterfeiter -o ./mocks/identifier.go . Identifier
type Identifier interface {
	ID() rbac.RoleID
}

// Provider of payload data of a role object
type Provider interface {
	Data() *data
}

// Incomplete role object
type Incomplete interface {
	Data() *data
}

// Complete role object
type Complete interface {
	Identifier
	Incomplete
}

type data struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Rules       rbac.RoleRules `json:"rules"`
}

// NewIncomplete creates a new incomplete role object
func NewIncomplete(title, description string, rules rbac.RoleRules) Incomplete {
	return &data{
		title,
		description,
		rules,
	}
}

func (d *data) Data() *data {
	return d
}

func (d *data) SetTitle(to string) *data {
	d.Title = to
	return d
}

func (d *data) SetDescription(to string) *data {
	d.Description = to
	return d
}

func (d *data) SetRules(to rbac.RoleRules) *data {
	d.Rules = to
	return d
}
