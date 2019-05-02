package role_test

import (
	"testing"

	"github.com/51st-state/api/pkg/apis/role"
	"github.com/51st-state/api/pkg/rbac"
)

func TestNewIncomplete(t *testing.T) {
	inc := role.NewIncomplete("title", "description", rbac.RoleRules{})

	if inc.Data() == nil {
		t.Fatal("the data should not be null")
	}
}

func TestIncompleteSetTitle(t *testing.T) {
	inc := role.NewIncomplete("title", "description", rbac.RoleRules{})

	inc.Data().SetTitle("anotherTitle")

	if inc.Data().Title != "anotherTitle" {
		t.Fatal("the title was not changed")
	}
}

func TestIncompleteSetDescription(t *testing.T) {
	inc := role.NewIncomplete("title", "description", rbac.RoleRules{})

	inc.Data().SetDescription("anotherDescription")

	if inc.Data().Description != "anotherDescription" {
		t.Fatal("the description was not changed")
	}
}

func TestIncompleteSetRules(t *testing.T) {
	inc := role.NewIncomplete("title", "description", rbac.RoleRules{
		"testRule",
	})

	inc.Data().SetRules(rbac.RoleRules{
		"testRule1",
		"testRule2",
	})
	if len(inc.Data().Rules) != 2 {
		t.Fatal("the rules were not set")
	}
}
