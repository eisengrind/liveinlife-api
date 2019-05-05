package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/51st-state/api/pkg/apis/user"
)

type wcfRepository struct {
	database *sql.DB
}

// NewWCFRepository for fetching specific Woltlab Community Framwork user data
// This function returns, as defined in the user package, only the hashed
// password and the email of the wcf user
func NewWCFRepository(db *sql.DB) user.WCFRepository {
	return &wcfRepository{db}
}

func (r *wcfRepository) GetInfo(ctx context.Context, id user.WCFUserID) (*user.WCFUserInfo, error) {
	var info user.WCFUserInfo
	hash := make([]byte, 0)

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT userId,
        email,
        password
        FROM wcf1_user
        WHERE userId = ?`,
		id,
	).Scan(
		&info.UserID,
		&info.Email,
		&hash,
	); err != nil {
		return nil, err
	}

	info.Password = newCompletePassword(hash)

	return &info, nil
}

func (r *wcfRepository) GetInfoByEmail(ctx context.Context, wcfEmail string) (*user.WCFUserInfo, error) {
	var info user.WCFUserInfo
	hash := make([]byte, 0)

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT userId,
        email,
        password
        FROM wcf1_user
        WHERE email = ?`,
		wcfEmail,
	).Scan(
		&info.UserID,
		&info.Email,
		&hash,
	); err != nil {
		return nil, err
	}

	info.Password = newCompletePassword(hash)

	return &info, nil
}

func (r *wcfRepository) GetInfoByUsername(ctx context.Context, username string) (*user.WCFUserInfo, error) {
	var info user.WCFUserInfo
	hash := make([]byte, 0)

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT userId,
        email,
        password
        FROM wcf1_user
        WHERE username = ?`,
		username,
	).Scan(
		&info.UserID,
		&info.Email,
		&hash,
	); err != nil {
		return nil, err
	}

	info.Password = newCompletePassword(hash)

	return &info, nil
}

type completePassword struct {
	hash []byte
}

func newCompletePassword(hash []byte) user.CompletePassword {
	return &completePassword{hash}
}

func (p *completePassword) Hash() []byte {
	return p.hash
}
