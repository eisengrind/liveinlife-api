package user

import "context"

// Repository of user objects
//go:generate counterfeiter -o ./mocks/repository.go . Repository
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	GetByGameSerialHash(context.Context, string) (Complete, error)
	GetByWCFUserID(context.Context, WCFUserID) (Complete, error)
	Create(context.Context, Incomplete) (Complete, error)
	Update(context.Context, Complete) error
	Delete(context.Context, Identifier) error
}

// WCFUserID of the user existing in the Woltlab Community Framwork database
type WCFUserID uint64

// WCFUserInfo of an existing user in the WCF database
type WCFUserInfo struct {
	UserID   WCFUserID
	Email    string
	Password CompletePassword
}

// WCFRepository of the Woltlab Community Framwork database
//go:generate counterfeiter -o ./mocks/wcf_repository.go . WCFRepository
type WCFRepository interface {
	GetInfo(context.Context, WCFUserID) (*WCFUserInfo, error)
	GetInfoByEmail(context.Context, string) (*WCFUserInfo, error)
	GetInfoByUsername(context.Context, string) (*WCFUserInfo, error)
}
