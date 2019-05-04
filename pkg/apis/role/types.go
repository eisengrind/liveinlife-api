package role

import (
	"encoding/json"

	"github.com/51st-state/api/pkg/rbac"
)

// Identifier of a role object
//go:generate counterfeiter -o ./mocks/identifier.go . Identifier
type Identifier interface {
	ID() rbac.RoleID
}

type identifier struct {
	id rbac.RoleID
}

func (i *identifier) ID() rbac.RoleID {
	return i.id
}

func newIdentifier(id rbac.RoleID) Identifier {
	return &identifier{id}
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

type complete struct {
	Identifier
	Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          rbac.RoleID    `json:"id"`
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Rules       rbac.RoleRules `json:"rules"`
	}{
		c.ID(),
		c.Data().Title,
		c.Data().Description,
		c.Data().Rules,
	})
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
