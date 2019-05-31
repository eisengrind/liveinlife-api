package inventory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/51st-state/api/pkg/apis/inventory"
	"github.com/51st-state/api/pkg/apis/inventory/mocks"
)

func TestManagerGet(t *testing.T) {
	repo := &mocks.FakeRepository{}
	m := inventory.NewManager(repo)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")

	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the identifier is invalid")
	}

	id.GUIDReturns("test")

	if _, err := m.Get(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}

type fakeComplete struct {
	inventory.Identifier
	inventory.Incomplete
}

func TestManagerCreate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	m := inventory.NewManager(repo)

	inc := inventory.NewIncomplete([]*inventory.Item{
		&inventory.Item{
			ID:     "",
			Amount: 0,
			Subset: -1.1,
		},
	})
	if _, err := m.Create(context.Background(), inc); err == nil {
		t.Fatal("the item name is invalid")
	}

	inc = inventory.NewIncomplete([]*inventory.Item{
		&inventory.Item{
			ID:     "testName",
			Amount: 0,
			Subset: -1.1,
		},
	})
	if _, err := m.Create(context.Background(), inc); err == nil {
		t.Fatal("the item amount is invalid")
	}

	inc = inventory.NewIncomplete([]*inventory.Item{
		&inventory.Item{
			ID:     "testName",
			Amount: 1,
			Subset: -1.1,
		},
	})
	if _, err := m.Create(context.Background(), inc); err == nil {
		t.Fatal("the item subset is invalid")
	}

	inc = inventory.NewIncomplete([]*inventory.Item{
		&inventory.Item{
			ID:     "testName",
			Amount: 1,
			Subset: -1,
		},
	})
	repo.CreateReturns(nil, errors.New("fake error"))
	if _, err := m.Create(context.Background(), inc); err == nil {
		t.Fatal("the create repository returns an error")
	}

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("testName")
	repo.CreateReturns(&fakeComplete{
		id,
		inc,
	}, nil)
	repo.AddItemReturns(errors.New("fake error"))
	if _, err := m.Create(context.Background(), inc); err == nil {
		t.Fatal("an item could not be added")
	}

	repo.AddItemReturns(nil)
	if _, err := m.Create(context.Background(), inc); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerAddItem(t *testing.T) {
	repo := &mocks.FakeRepository{}
	m := inventory.NewManager(repo)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")

	if err := m.AddItem(context.Background(), id, nil); err == nil {
		t.Fatal("the id is invalid")
	}

	id.GUIDReturns("testName")
	item := &inventory.Item{
		ID:     "",
		Amount: 0,
		Subset: -1.1,
	}

	if err := m.AddItem(context.Background(), id, item); err == nil {
		t.Fatal("the item id is invalid")
	}

	item.ID = "testName"
	if err := m.AddItem(context.Background(), id, item); err == nil {
		t.Fatal("the item amount is invalid")
	}

	item.Amount = 1
	if err := m.AddItem(context.Background(), id, item); err == nil {
		t.Fatal("the item subset is invalid")
	}

	item.Subset = -1
	if err := m.AddItem(context.Background(), id, item); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerRemoveItem(t *testing.T) {
	repo := &mocks.FakeRepository{}
	m := inventory.NewManager(repo)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")
	item := &inventory.Item{
		ID:     "",
		Amount: 0,
		Subset: -1.1,
	}
	if err := m.RemoveItem(context.Background(), id, item); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("testName")
	if err := m.RemoveItem(context.Background(), id, item); err == nil {
		t.Fatal("the item id is invalid")
	}

	item.ID = "testName"
	if err := m.RemoveItem(context.Background(), id, item); err == nil {
		t.Fatal("the item amount is invalid")
	}

	item.Amount = 1
	if err := m.RemoveItem(context.Background(), id, item); err == nil {
		t.Fatal("the item subset is invalid")
	}

	item.Subset = -1
	if err := m.RemoveItem(context.Background(), id, item); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerDelete(t *testing.T) {
	repo := &mocks.FakeRepository{}
	m := inventory.NewManager(repo)

	id := &mocks.FakeIdentifier{}
	id.GUIDReturns("")
	if err := m.Delete(context.Background(), id); err == nil {
		t.Fatal("the guid is invalid")
	}

	id.GUIDReturns("testName")
	if err := m.Delete(context.Background(), id); err != nil {
		t.Fatal("there should be no error")
	}
}
