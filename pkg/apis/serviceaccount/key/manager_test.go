package key_test

import (
	"context"
	"errors"
	"testing"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"
	"github.com/51st-state/api/pkg/apis/serviceaccount/key/mocks"
	serviceaccountMocks "github.com/51st-state/api/pkg/apis/serviceaccount/mocks"
)

func TestManagerGet(t *testing.T) {
	repo := &mocks.FakeRepository{}
	saMgr := &serviceaccountMocks.FakeManager{}
	manager := key.NewManager(repo, saMgr)

	id := &mocks.FakeIdentifier{}

	if _, err := manager.Get(context.Background(), id); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("test")

	if _, err := manager.Get(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

type fakeComplete struct {
	key.Identifier
	key.Incomplete
}

func TestManagerUpdate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	saMgr := &serviceaccountMocks.FakeManager{}
	manager := key.NewManager(repo, saMgr)

	id := &mocks.FakeIdentifier{}
	inc := key.NewIncomplete("", "")
	c := &fakeComplete{
		id,
		inc,
	}

	if err := manager.Update(context.Background(), c); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("test")

	if err := manager.Update(context.Background(), c); err == nil {
		t.Fatal("the name is invalid")
	}

	inc.Data().Name = "testName"

	if err := manager.Update(context.Background(), c); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerDelete(t *testing.T) {
	repo := &mocks.FakeRepository{}
	saMgr := &serviceaccountMocks.FakeManager{}
	manager := key.NewManager(repo, saMgr)

	id := &mocks.FakeIdentifier{}

	if err := manager.Delete(context.Background(), id); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("test")
	saMgr.GetReturns(nil, nil)

	if err := manager.Delete(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerCreate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	saMgr := &serviceaccountMocks.FakeManager{}
	manager := key.NewManager(repo, saMgr)

	sid := &serviceaccountMocks.FakeIdentifier{}
	sid.GUIDReturns("")
	saMgr.GetReturns(nil, errors.New("fake error"))
	inc := key.NewIncomplete("", "")
	inc.Data().ServiceAccountGUID = ""

	if _, err := manager.Create(context.Background(), inc); err == nil {
		t.Fatal("the service account does not exist")
	}

	inc.Data().ServiceAccountGUID = "test"
	saMgr.GetReturns(nil, nil)

	if _, err := manager.Create(context.Background(), inc); err == nil {
		t.Fatal("the name is invalid")
	}

	inc.Data().Name = "test"
	repo.CreateReturns(nil, errors.New("fake error"))

	if _, err := manager.Create(context.Background(), inc); err == nil {
		t.Fatal("the repository method Create() returns an error")
	}

	repo.CreateReturns(&fakeComplete{
		key.NewIdentifier("test"),
		key.NewIncomplete("test", "test"),
	}, nil)

	if _, err := manager.Create(context.Background(), inc); err != nil {
		t.Fatal("there should be no error")
	}
}
