package role_test

import (
	"context"
	"errors"
	"testing"

	"github.com/51st-state/api/pkg/apis/role"
	"github.com/51st-state/api/pkg/apis/role/mocks"
	"github.com/51st-state/api/pkg/rbac"

	rbacMocks "github.com/51st-state/api/pkg/rbac/mocks"
)

type fakeComplete struct {
	role.Identifier
	role.Incomplete
}

func TestManagerGet(t *testing.T) {
	control := &rbacMocks.FakeControl{}
	repo := &mocks.FakeRepository{}
	m := role.NewManager(repo, control)

	id := &mocks.FakeIdentifier{}
	id.IDReturns("")
	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the id of the role is empty")
	}

	id.IDReturns("testid")

	repo.GetReturns(nil, errors.New("fake error"))
	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.GetReturns(&fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}, nil)
	control.GetRoleRulesReturns(nil, errors.New("fake error"))

	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the rbac service returns an error")
	}

	control.GetRoleRulesReturns(rbac.RoleRules{
		"testRule",
	}, nil)

	c, err := m.Get(context.Background(), id)
	if err != nil {
		t.Fatal("there should be no error")
	}

	if c.Data().Rules[0] != "testRule" {
		t.Fatal("the returned rules are not equal")
	}
}

func TestManagerSet(t *testing.T) {
	control := &rbacMocks.FakeControl{}
	repo := &mocks.FakeRepository{}
	m := role.NewManager(repo, control)

	id := &mocks.FakeIdentifier{}
	id.IDReturns("")
	if err := m.Set(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the id of the role is empty")
	}

	id.IDReturns("testid")
	repo.UpdateReturns(errors.New("fake error"))
	if err := m.Set(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.UpdateReturns(nil)
	control.SetRoleRulesReturns(errors.New("fake error"))
	if err := m.Set(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the rbac repository returns an error")
	}

	control.SetRoleRulesReturns(nil)
	if err := m.Set(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerCreate(t *testing.T) {
	control := &rbacMocks.FakeControl{}
	repo := &mocks.FakeRepository{}
	m := role.NewManager(repo, control)

	id := &mocks.FakeIdentifier{}
	id.IDReturns("")
	if err := m.Create(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the id of the role is empty")
	}

	id.IDReturns("testid")
	repo.CreateReturns(errors.New("fake error"))
	if err := m.Create(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.CreateReturns(nil)
	control.SetRoleRulesReturns(errors.New("fake error"))
	if err := m.Create(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err == nil {
		t.Fatal("the rbac control returns an error")
	}

	control.SetRoleRulesReturns(nil)

	if err := m.Create(context.Background(), &fakeComplete{
		id,
		role.NewIncomplete("title", "description", rbac.RoleRules{}),
	}); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerDelete(t *testing.T) {
	control := &rbacMocks.FakeControl{}
	repo := &mocks.FakeRepository{}
	m := role.NewManager(repo, control)

	id := &mocks.FakeIdentifier{}
	id.IDReturns("")
	if err := m.Delete(context.Background(), id); err == nil {
		t.Fatal("the id of the role is empty")
	}

	id.IDReturns("testid")

	control.SetRoleRulesReturns(errors.New("fake error"))
	if err := m.Delete(context.Background(), id); err == nil {
		t.Fatal("the rbac control returns an error")
	}

	control.SetRoleRulesReturns(nil)
	repo.DeleteReturns(errors.New("fake error"))

	if err := m.Delete(context.Background(), id); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.DeleteReturns(nil)
	if err := m.Delete(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}
