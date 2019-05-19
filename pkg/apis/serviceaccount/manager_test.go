package serviceaccount_test

import (
	"context"
	"errors"
	"testing"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
	"github.com/51st-state/api/pkg/apis/serviceaccount/mocks"
	"github.com/51st-state/api/pkg/rbac"
	rbacMocks "github.com/51st-state/api/pkg/rbac/mocks"
)

func TestManagerGet(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")

	if _, err := manager.Get(context.Background(), id); err == nil {
		t.Fatal("there has to be an error since the guid is invalid")
	}

	id.GUIDReturns("test")

	if _, err := manager.Get(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

type fakeComplete struct {
	serviceaccount.Identifier
	serviceaccount.Incomplete
}

func TestManagerUpdate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")

	inc := serviceaccount.NewIncomplete("", "")

	complete := &fakeComplete{id, inc}

	if err := manager.Update(context.Background(), complete); err == nil {
		t.Fatal("the given guid is invalid")
	}

	id.GUIDReturns("test")

	if err := manager.Update(context.Background(), complete); err == nil {
		t.Fatal("the name is invalid")
	}

	inc.Data().Name = "test"

	if err := manager.Update(context.Background(), complete); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerCreate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	inc := serviceaccount.NewIncomplete("", "")

	if _, err := manager.Create(context.Background(), inc); err == nil {
		t.Fatal("there has to be an error since the name is invalid")
	}

	inc.Data().Name = "test"

	if _, err := manager.Create(context.Background(), inc); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerDelete(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")

	if err := manager.Delete(context.Background(), id); err == nil {
		t.Fatal("there has to be an error since the guid is invalid")
	}

	id.GUIDReturns("test")

	if err := manager.Delete(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerGetRoles(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	id := &mocks.FakeIdentifier{}

	if _, err := manager.GetRoles(context.Background(), id); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("test")
	rbControl.GetAccountRolesReturns(nil, errors.New("fake error"))
	if _, err := manager.GetRoles(context.Background(), id); err == nil {
		t.Fatal("rbac returns an error")
	}

	rbControl.GetAccountRolesReturns(rbac.AccountRoles{}, nil)
	if _, err := manager.GetRoles(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerSetRoles(t *testing.T) {
	repo := &mocks.FakeRepository{}
	rbControl := &rbacMocks.FakeControl{}
	manager := serviceaccount.NewManager(repo, rbControl)

	id := &mocks.FakeIdentifier{}

	if err := manager.SetRoles(context.Background(), id, rbac.AccountRoles{}); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("test")

	if err := manager.SetRoles(context.Background(), id, rbac.AccountRoles{}); err != nil {
		t.Fatal("there should be no error")
	}
}
