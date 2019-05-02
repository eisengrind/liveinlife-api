package user

import (
	"context"
	"fmt"
	"regexp"

	"github.com/51st-state/api/pkg/event"
	"github.com/51st-state/api/pkg/rbac"

	"github.com/pkg/errors"

	"github.com/51st-state/api/pkg/bcrypt"
)

// Manager of user objects
type Manager struct {
	repository    Repository
	wcfRepository WCFRepository
	event         *event.Producer
	rbac          rbac.Control
}

// NewManager for user objects
//go:generate protoc -I./../../../../../../  -I ./proto --go_out=plugins=grpc:./proto ./proto/manager.proto
func NewManager(r Repository, wcf WCFRepository, prod *event.Producer, rb rbac.Control) *Manager {
	return &Manager{
		r,
		wcf,
		prod,
		rb,
	}
}

var errInvalidUUID = errors.New("invalid user uuid given")

// Get an user object
func (m *Manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if id.UUID() == "" {
		return nil, errInvalidUUID
	}

	return m.repository.Get(ctx, id)
}

// GetByGameSerialHash returns a user filtered by its unique game serial hash
func (m *Manager) GetByGameSerialHash(ctx context.Context, hash string) (Complete, error) {
	if hash == "" {
		return nil, errInvalidGameSerialHash
	}

	return m.repository.GetByGameSerialHash(ctx, hash)
}

// GetByWCFUserID returns an user filtered by its wcf user id
func (m *Manager) GetByWCFUserID(ctx context.Context, wcfUserID WCFUserID) (Complete, error) {
	if wcfUserID == 0 {
		return nil, errInvalidWCFUserID
	}

	return m.repository.GetByWCFUserID(ctx, wcfUserID)
}

var errInvalidWCFUserID = errors.New("invalid woltlab community framework user id")

// Create an user object
func (m *Manager) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	if inc.Data().WCFUserID == 0 {
		return nil, errInvalidWCFUserID
	}

	if _, err := m.wcfRepository.GetInfo(ctx, inc.Data().WCFUserID); err != nil {
		return nil, err
	}

	c, err := m.repository.Create(ctx, inc)
	if err != nil {
		return nil, err
	}

	return c, m.event.Produce(ctx, CreatedEventID, &CreatedEvent{
		&event.PayloadMeta{
			Version: "1",
		},
		c,
	})
}

var (
	errInvalidEmailFormat = errors.New("invalid email format")
	emailRegexp           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func (m *Manager) checkEmail(ctx context.Context, email string) error {
	if !emailRegexp.MatchString(email) {
		return errInvalidEmailFormat
	}

	return nil
}

// Delete an user object
func (m *Manager) Delete(ctx context.Context, id Identifier) error {
	if id.UUID() == "" {
		return errInvalidUUID
	}

	if err := m.repository.Delete(ctx, id); err != nil {
		return err
	}

	return m.event.Produce(ctx, DeletedEventID, &DeletedEvent{
		&event.PayloadMeta{
			Version: "1",
		},
		id,
	})
}

// GetWCFInfoByEmail of a wcf user
func (m *Manager) GetWCFInfoByEmail(ctx context.Context, email string) (*WCFUserInfo, error) {
	return m.wcfRepository.GetInfoByEmail(ctx, email)
}

// GetWCFInfoByUsername of a wcf user
func (m *Manager) GetWCFInfoByUsername(ctx context.Context, username string) (*WCFUserInfo, error) {
	return m.wcfRepository.GetInfoByUsername(ctx, username)
}

var errInvalidGameSerialHash = errors.New("invalid game serial hash")

// Update an user objects data
func (m *Manager) Update(ctx context.Context, c Complete) error {
	if c.UUID() == "" {
		return errInvalidUUID
	}

	if c.Data().WCFUserID == 0 {
		return errInvalidWCFUserID
	}

	if _, err := m.wcfRepository.GetInfo(ctx, c.Data().WCFUserID); err != nil {
		return err
	}

	if err := m.repository.Update(ctx, c); err != nil {
		return err
	}

	return m.event.Produce(ctx, UpdatedEventID, &UpdatedEvent{
		&event.PayloadMeta{
			Version: "1",
		},
		c,
	})
}

// CheckPassword of a user
func (m *Manager) CheckPassword(ctx context.Context, id Identifier, incPw IncompletePassword) error {
	if id.UUID() == "" {
		return errInvalidUUID
	}

	compl, err := m.Get(ctx, id)
	if err != nil {
		return err
	}

	wcfInfo, err := m.wcfRepository.GetInfo(ctx, compl.Data().WCFUserID)
	if err != nil {
		return err
	}

	pwFirstHash, err := getFirstPasswordHash(wcfInfo.Password.Hash(), []byte(incPw.Password()))
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(wcfInfo.Password.Hash(), pwFirstHash)
}

const (
	majorBcryptVersion = '2'
	minorBcryptVersion = 'a'
)

func getFirstPasswordHash(hash, password []byte) ([]byte, error) {
	hashInfo, err := bcrypt.NewFromHash(hash)
	if err != nil {
		return nil, err
	}

	pwHashOnly, err := bcrypt.Bcrypt(password, hashInfo.Cost, hashInfo.Salt)
	if err != nil {
		return nil, err
	}

	return bcrypt.NewHash(
		pwHashOnly,
		hashInfo.Salt,
		hashInfo.Cost,
		majorBcryptVersion,
		minorBcryptVersion,
	).Hash(), nil
}

// GetRoles of a user
func (m *Manager) GetRoles(ctx context.Context, id Identifier) (rbac.SubjectRoles, error) {
	roles, err := m.rbac.GetSubjectRoles(ctx, rbac.SubjectID(fmt.Sprintf(
		"user/%s",
		id.UUID(),
	)))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// SetRoles of a user
func (m *Manager) SetRoles(ctx context.Context, id Identifier, roles rbac.SubjectRoles) error {
	return m.rbac.SetSubjectRoles(ctx, rbac.SubjectID(fmt.Sprintf(
		"user/%s",
		id.UUID(),
	)), roles)
}
