package rbac_test

import (
	"testing"

	"github.com/51st-state/api/pkg/rbac"
)

func TestRoleRulesContains(t *testing.T) {
	rules := rbac.RoleRules{
		"testRule1",
		"testRule2",
	}

	if !rules.Contains("testRule1") {
		t.Fatal("this rule exists in given role rules")
	}

	if rules.Contains("test") {
		t.Fatal("this rule does not exist")
	}
}
