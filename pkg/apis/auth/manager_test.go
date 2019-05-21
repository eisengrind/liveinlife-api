package auth_test

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"

	"github.com/51st-state/api/pkg/token"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/51st-state/api/pkg/apis/auth"
	"github.com/51st-state/api/pkg/apis/user"

	"github.com/51st-state/api/pkg/keys"

	mocks "github.com/51st-state/api/pkg/apis/auth/mocks"
	keyMocks "github.com/51st-state/api/pkg/apis/serviceaccount/key/mocks"
	userMocks "github.com/51st-state/api/pkg/apis/user/mocks"
	"github.com/51st-state/api/test"
)

var testPrivateKey *rsa.PrivateKey

func init() {
	var err error
	testPrivateKey, err = keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		panic(err)
	}
}

func TestNewManager(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	auth.NewManager(prvKey, repo, userManager, nil, nil)
}

type testCredentials struct {
	name     string
	password string
}

func (t *testCredentials) Name() string {
	return t.name
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

const privateKeyTestResult = `{"service_account_guid":"","guid":"","private_key":"-----BEGIN PRIVATE KEY-----\n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\n-----END PRIVATE KEY-----\n"}`
const invalidPrivateKeyTestResult = `{"service_account_guid":"","guid":"","private_key":"-----BEGIN PRIVATE KEY-----\n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\n-----END PRIVATE KEY-----\n"}`

type fakeCompleteServiceAccountKey struct {
	key.Identifier
	key.Incomplete
}

func TestManagerLogin(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}
	keySA := &keyMocks.FakeManager{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}
	manager := auth.NewManager(prvKey, repo, userManager, nil, keySA)

	id := &userMocks.FakeIdentifier{}
	id.UUIDReturns("user/test")
	userManager.GetWCFInfoReturns(&user.WCFUserInfo{
		UserID: 1,
	}, nil)
	userManager.GetByWCFUserIDReturns(newComplete(
		id,
		user.NewIncomplete(1, "", "", "", false),
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
		user.NewIncomplete(1, "", "", "", false),
	), nil)
	userManager.CheckPasswordReturns(errors.New("fake error"))

	if _, err := manager.Login(context.Background(), &testCredentials{
		"user/test",
		"1234",
	}); err == nil {
		t.Fatal("create user returns an error")
	}

	c := &testCredentials{
		"_json_key",
		privateKeyTestResult + ".",
	}
	if _, err := manager.Login(context.Background(), c); err == nil {
		t.Fatal("the json key has an invalid format")
	}

	c.password = privateKeyTestResult
	keySA.GetReturns(nil, errors.New("fake error"))
	if _, err := manager.Login(context.Background(), c); err == nil {
		t.Fatal("the service account key does not exist")
	}

	inc := key.NewIncomplete("test", "test")
	inc.Data().PublicKey = &testPrivateKey.PublicKey
	keySA.GetReturns(&fakeCompleteServiceAccountKey{
		key.NewIdentifier("test"),
		inc,
	}, nil)

	if _, err := manager.Login(context.Background(), c); err != nil {
		t.Fatal("there should be no error")
	}
}

func TestManagerRefreshToken(t *testing.T) {
	userManager := &userMocks.FakeManager{}
	repo := &mocks.FakeRepository{}

	prvKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}
	manager := auth.NewManager(prvKey, repo, userManager, nil, nil)

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
