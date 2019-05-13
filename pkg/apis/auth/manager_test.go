package auth_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/51st-state/api/pkg/token"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/51st-state/api/pkg/apis/auth"
	"github.com/51st-state/api/pkg/apis/user"

	"github.com/51st-state/api/pkg/keys"

	mocks "github.com/51st-state/api/pkg/apis/auth/mocks"
	userMocks "github.com/51st-state/api/pkg/apis/user/mocks"
	"github.com/51st-state/api/test"
)

func TestNewManager(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	auth.NewManager(prvKey, repo, userManager, nil)
}

type testCredentials struct {
	username string
	password string
}

func (t *testCredentials) Username() string {
	return t.username
}

func (t *testCredentials) Password() string {
	return t.password
}

type fakeComplete struct {
	user.Identifier
	user.Incomplete
}

func newComplete(id user.Identifier, inc user.Incomplete) user.Complete {
	return &fakeComplete{id, inc}
}

func TestManagerLogin(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}
	manager := auth.NewManager(prvKey, repo, userManager, nil)

	id := &userMocks.FakeIdentifier{}
	id.UUIDReturns("user/test")
	userManager.GetWCFInfoReturns(&user.WCFUserInfo{
		UserID: 1,
	}, nil)
	userManager.GetByWCFUserIDReturns(newComplete(
		id,
		user.NewIncomplete(1, "", false),
	), nil)
	userManager.CheckPasswordReturns(nil)

	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err != nil {
		t.Fatal("all parameters are correct")
	}

	repo.LoginAttemptsCountSinceReturns(0, errors.New("fake error"))
	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("the repository returns an error")
	}

	repo.LoginAttemptsCountSinceReturns(1, nil)
	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("the count of login attempts is = 1")
	}

	userManager.GetWCFInfoReturns(nil, errors.New("fake error"))
	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("the wcf repo returns an error")
	}

	userManager.GetWCFInfoReturns(&user.WCFUserInfo{
		UserID: 1,
	}, nil)
	userManager.GetByWCFUserIDReturns(nil, sql.ErrNoRows)
	userManager.CreateReturns(nil, errors.New("fake error"))
	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("create user returns an error")
	}

	userManager.GetByWCFUserIDReturns(nil, sql.ErrConnDone)
	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("the wcf user id returns an error")
	}

	userManager.GetByWCFUserIDReturns(newComplete(
		id,
		user.NewIncomplete(1, "", false),
	), nil)
	userManager.CheckPasswordReturns(errors.New("fake error"))

	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("create user returns an error")
	}
}

func TestManagerRefreshToken(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}
	manager := auth.NewManager(prvKey, repo, userManager, nil)

	if _, err := manager.RefreshToken(context.Background(), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "defaultFake",
	}, &token.User{}), nil); err == nil {
		t.Fatal("the given audience for the access token is invalid")
	}

	if _, err := manager.RefreshToken(context.Background(), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "default",
	}, &token.User{}), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "auth/refreshFake",
	}, &token.User{})); err == nil {
		t.Fatal("the given audience for the refresh token is invalid")
	}

	if _, err := manager.RefreshToken(context.Background(), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "default",
	}, &token.User{
		ID:   "1234",
		Type: "user",
	}), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "auth/refresh",
	}, &token.User{
		ID:   "1234",
		Type: "service_account",
	})); err == nil {
		t.Fatal("the user ids are invalid")
	}

	if _, err := manager.RefreshToken(context.Background(), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "default",
	}, &token.User{
		ID:   "1234",
		Type: "user",
	}), token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Audience:  "auth/refresh",
	}, &token.User{
		ID:   "1234",
		Type: "user",
	})); err != nil {
		t.Fatal("there should be no error")
	}
}
