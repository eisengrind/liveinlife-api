package rbac_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/51st-state/api/pkg/rbac"
	"github.com/51st-state/api/pkg/rbac/mocks"
)

func TestNewControl(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)
	fmt.Println(ctrl)
}

func TestControlGetRoleRules(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)

	if _, err := ctrl.GetRoleRules(context.Background(), ""); err == nil {
		t.Fatal("empty role id")
	}

	if _, err := ctrl.GetRoleRules(context.Background(), "testid"); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestControlSetRoleRules(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)

	if err := ctrl.SetRoleRules(context.Background(), "", rbac.RoleRules{}); err == nil {
		t.Fatal("empty role id")
	}

	if err := ctrl.SetRoleRules(context.Background(), "testid", rbac.RoleRules{
		"",
	}); err == nil {
		t.Fatal("empty rule id")
	}

	if err := ctrl.SetRoleRules(context.Background(), "testid", rbac.RoleRules{}); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestControlGetAccountRoles(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)

	if _, err := ctrl.GetAccountRoles(context.Background(), ""); err == nil {
		t.Fatal("empty role id")
	}

	if _, err := ctrl.GetAccountRoles(context.Background(), "testid"); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestControlSetAccountRoles(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)

	if err := ctrl.SetAccountRoles(context.Background(), "", rbac.AccountRoles{}); err == nil {
		t.Fatal("empty account id")
	}

	if err := ctrl.SetAccountRoles(context.Background(), "testid", rbac.AccountRoles{
		"",
	}); err == nil {
		t.Fatal("empty role id")
	}

	if err := ctrl.SetAccountRoles(context.Background(), "testid", rbac.AccountRoles{}); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestControlIsAccountAllowed(t *testing.T) {
	repo := &mocks.FakeRepository{}
	ctrl := rbac.NewControl(repo)

	if err := ctrl.IsAccountAllowed(context.Background(), "", "ruleID"); err == nil {
		t.Fatal("empty account id")
	}

	if err := ctrl.IsAccountAllowed(context.Background(), "accountID", ""); err == nil {
		t.Fatal("empty rule")
	}

	repo.GetAccountRuleCountReturns(0, errors.New("fake error"))
	if err := ctrl.IsAccountAllowed(context.Background(), "accountID", "rule"); err == nil {
		t.Fatal("repository returns an error")
	}

	repo.GetAccountRuleCountReturns(0, nil)
	if err := ctrl.IsAccountAllowed(context.Background(), "accountID", "rule"); err == nil {
		t.Fatal("the account does have no permission")
	}

	repo.GetAccountRuleCountReturns(1, nil)
	if err := ctrl.IsAccountAllowed(context.Background(), "accountID", "rule"); err != nil {
		t.Fatal("there should be no error")
	}
}
