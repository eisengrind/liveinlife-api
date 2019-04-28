package user_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/51st-state/api/pkg/event"

	"github.com/51st-state/api/pkg/apis/user"
	"github.com/51st-state/api/pkg/apis/user/mocks"
	pubsubMocks "github.com/51st-state/api/pkg/pubsub/mocks"
	"github.com/pkg/errors"
)

type fakeComplete struct {
	user.Identifier
	user.Incomplete
}

func (c *fakeComplete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{}{})
}

func TestManagerGet(t *testing.T) {
	repo := &mocks.FakeRepository{}
	wcfRepo := &mocks.FakeWCFRepository{}

	m := user.NewManager(repo, wcfRepo, event.NewProducer(&pubsubMocks.FakeProducer{}))

	fComplPassword := &mocks.FakeCompletePassword{}
	fComplPassword.HashReturns([]byte("testHash"))

	wcfRepo.GetInfoReturns(&user.WCFUserInfo{
		Email:    "test@test.com",
		Password: fComplPassword,
	}, nil)
	repo.GetReturns(nil, errors.New("fake error"))

	id := &mocks.FakeIdentifier{}
	id.UUIDReturns("test")

	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the function has an error given")
	}

	fakeIdent := &mocks.FakeIdentifier{}
	fakeIdent.UUIDReturns("test")

	repo.GetReturns(&fakeComplete{
		fakeIdent,
		user.NewIncomplete(1, "testSerialHash", false),
	}, nil)
	if _, err := m.Get(context.Background(), id); err != nil {
		t.Fatal("there is no error")
	}

	id.UUIDReturns("")
	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("Get() does not accept empty uuids")
	}

	wcfRepo.GetInfoReturns(nil, errors.New("fake error"))
	if _, err := m.Get(context.Background(), id); err == nil {
		t.Fatal("the wcf repository returns an error")
	}
}

func TestManagerCreate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	wcfRepo := &mocks.FakeWCFRepository{}

	m := user.NewManager(repo, wcfRepo, event.NewProducer(&pubsubMocks.FakeProducer{}))

	validIncomplete := user.NewIncomplete(1, "exampleHash", false)
	if _, err := m.Create(
		context.Background(),
		validIncomplete,
	); err != nil {
		t.Fatal("the given parameters are both correct")
	}

	if _, err := m.Create(
		context.Background(),
		user.NewIncomplete(0, "exampleHash", false),
	); err == nil {
		t.Fatal("the wcf user id is invalid")
	}

	wcfRepo.GetInfoReturns(nil, errors.New("fake error"))
	if _, err := m.Create(
		context.Background(),
		validIncomplete,
	); err == nil {
		t.Fatal("the wcf repository returns an error")
	}

	repo.GetReturns(nil, errors.New("fake error"))
	if _, err := m.Create(
		context.Background(),
		validIncomplete,
	); err == nil {
		t.Fatal("the repository returns an error")
	}
}

type fakeIdentifier struct {
	UUIDStr string `json:"id_str"`
}

func (id *fakeIdentifier) UUID() string {
	return id.UUIDStr
}

func TestManagerDelete(t *testing.T) {
	repo := &mocks.FakeRepository{}
	wcfRepo := &mocks.FakeWCFRepository{}

	m := user.NewManager(repo, wcfRepo, event.NewProducer(&pubsubMocks.FakeProducer{}))

	id := &fakeIdentifier{""}

	if err := m.Delete(context.Background(), id); err == nil {
		t.Fatal("the given uuid is empty")
	}

	id.UUIDStr = "test"
	if err := m.Delete(context.Background(), id); err != nil {
		t.Fatal("given request is correct")
	}
}

func TestManagerUpdate(t *testing.T) {
	repo := &mocks.FakeRepository{}
	wcfRepo := &mocks.FakeWCFRepository{}

	m := user.NewManager(repo, wcfRepo, event.NewProducer(&pubsubMocks.FakeProducer{}))

	id := &mocks.FakeIdentifier{}
	id.UUIDReturns("")

	if err := m.Update(context.Background(), &fakeComplete{
		id,
		user.NewIncomplete(1, "exampleHash", false),
	}); err == nil {
		t.Fatal("the given uuid is empty")
	}

	id.UUIDReturns("test")

	if err := m.Update(
		context.Background(),
		&fakeComplete{
			id,
			user.NewIncomplete(1, "exampleHash", false),
		},
	); err != nil {
		t.Fatal("the given parameters are both correct")
	}

	wcfRepo.GetInfoReturns(nil, errors.New("fake error"))
	if err := m.Update(
		context.Background(),
		&fakeComplete{
			id,
			user.NewIncomplete(1, "exampleGameSerialHash", false),
		},
	); err == nil {
		t.Fatal("the wcf repository returns an error")
	}

	repo.UpdateReturns(errors.New("fake error"))
	if err := m.Update(
		context.Background(),
		&fakeComplete{
			id,
			user.NewIncomplete(1, "exampleGameSerialHash", false),
		},
	); err == nil {
		t.Fatal("the repository returns an error")
	}
}

func TestManagerCheckPassword(t *testing.T) {
	repo := &mocks.FakeRepository{}
	wcfRepo := &mocks.FakeWCFRepository{}

	m := user.NewManager(repo, wcfRepo, event.NewProducer(&pubsubMocks.FakeProducer{}))

	id := &mocks.FakeIdentifier{}
	id.UUIDReturns("")

	if err := m.CheckPassword(context.Background(), id, &mocks.FakeIncompletePassword{}); err == nil {
		t.Fatal("the uuid is invalid")
	}

	id.UUIDReturns("test")
	repo.GetReturns(nil, errors.New("fake error"))

	if err := m.CheckPassword(context.Background(), id, &mocks.FakeIncompletePassword{}); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.GetReturns(&fakeComplete{
		id,
		user.NewIncomplete(1, "testSerialHash", false),
	}, nil)
	wcfRepo.GetInfoReturns(nil, errors.New("fake error"))

	if err := m.CheckPassword(context.Background(), id, &mocks.FakeIncompletePassword{}); err == nil {
		t.Fatal("the wcf repository returns an error")
	}

	fakePw := &mocks.FakeCompletePassword{}
	fakePw.HashReturns([]byte(""))

	wcfRepo.GetInfoReturns(&user.WCFUserInfo{
		Email:    "test@example.com",
		Password: fakePw,
	}, nil)

	if err := m.CheckPassword(context.Background(), id, &mocks.FakeIncompletePassword{}); err == nil {
		t.Fatal("the first password hash generation returned an error")
	}

	fakePw.HashReturns([]byte("$2a$08$asdqweasdyxcasdqweasdOr9sGCG7KJN.58c3i4IIICPOMS9uUp9S"))

	fakeInc := &mocks.FakeIncompletePassword{}
	fakeInc.PasswordReturns("root")

	if err := m.CheckPassword(context.Background(), id, fakeInc); err != nil {
		t.Fatal("the password is not equal")
	}

	fakeInc.PasswordReturns("roo")

	if err := m.CheckPassword(context.Background(), id, fakeInc); err == nil {
		t.Fatal("the password is definetely invalid")
	}
}