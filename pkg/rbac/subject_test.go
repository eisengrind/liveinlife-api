package rbac_test

import (
	"testing"

	"github.com/51st-state/api/pkg/rbac"
)

func TestSubjectRolesContains(t *testing.T) {
	roles := rbac.SubjectRoles{
		"testRole1",
		"testRole2",
	}

	if !roles.Contains("testRole1") {
		t.Fatal("this role exists in given subject roles")
	}

	if roles.Contains("test") {
		t.Fatal("this role does not exist")
	}
}
