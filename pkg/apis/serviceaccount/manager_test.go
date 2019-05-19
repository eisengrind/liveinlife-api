package serviceaccount_test

import (
	"context"
	"testing"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
	"github.com/51st-state/api/pkg/apis/serviceaccount/mocks"
)

func TestManagerGet(t *testing.T) {
	repo := &mocks.FakeRepository{}
	manager := serviceaccount.NewManager(repo)

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
	manager := serviceaccount.NewManager(repo)

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
	manager := serviceaccount.NewManager(repo)

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
	manager := serviceaccount.NewManager(repo)

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
