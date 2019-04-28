package user

// Identifier of a user object
//go:generate counterfeiter -o ./mocks/identifier.go . Identifier
type Identifier interface {
	UUID() string
}

type identifier struct {
	uuid string
}

func newIdentifier(uuid string) Identifier {
	return &identifier{
		uuid,
	}
}

func (i *identifier) UUID() string {
	return i.uuid
}

// Provider of user data
type Provider interface {
	Data() *data
}

// Incomplete user object
type Incomplete interface {
	Provider
}

// IncompletePassword of a user (object).
// A complete password is containing a password and a salt (see CompletePassword)
//go:generate counterfeiter -o ./mocks/incomplete_password.go . IncompletePassword
type IncompletePassword interface {
	Password() string
}

type incompletePassword struct {
	pw string
}

// newIncompletePassword of a user
func newIncompletePassword(pw string) IncompletePassword {
	return &incompletePassword{
		pw,
	}
}

func (i *incompletePassword) Password() string {
	return i.pw
}

// CompletePassword of a user.
// In general, this password interface is returned only
// by the regarding repository and thus, is only created by it.
//go:generate counterfeiter -o ./mocks/complete_password.go . CompletePassword
type CompletePassword interface {
	Hash() []byte
}

type completePassword struct {
	hash []byte
}

// newCompletePassword of a user
func newCompletePassword(hash []byte) CompletePassword {
	return &completePassword{
		hash,
	}
}

func (c *completePassword) Hash() []byte {
	return c.hash
}

// PasswordChecker of user objects.
// Used to check

// Complete user object
type Complete interface {
	Identifier
	Incomplete
}

type complete struct {
	Identifier
	Incomplete
}

func newComplete(id Identifier, inc Incomplete) Complete {
	return &complete{
		id,
		inc,
	}
}

type data struct {
	WCFUserID      WCFUserID `json:"wcf_user_id"`
	GameSerialHash string    `json:"game_serial_hash"`
	Banned         bool      `json:"banned"`
}

// NewIncomplete user object
func NewIncomplete(wcfUserID WCFUserID, gameSerialHash string, status bool) Incomplete {
	return &data{wcfUserID, gameSerialHash, status}
}

func (d *data) Data() *data {
	return d
}

func (d *data) SetWCFUserID(to WCFUserID) *data {
	d.WCFUserID = to
	return d
}

func (d *data) SetGameSerialHash(to string) *data {
	d.GameSerialHash = to
	return d
}

func (d *data) SetBanned(to bool) *data {
	d.Banned = to
	return d
}
