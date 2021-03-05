package storage

import (
	"context"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"time"
)

type rootAccountModel struct {
	ID        uint64 `gorm:"primary_key"`
	UserID    uint64
	Type      string
	State     string
	Metadata  []byte
	CreatedAt int64
	UpdatedAt int64
}

func (a *rootAccountModel) TableName() string {
	return "root_accounts"
}

func (a *rootAccountModel) model(account domain.RootAccount) {
	a.ID = account.ID
	a.UserID = account.UserID
	a.Type = account.Type.String()
	a.State = account.State.String()
	a.Metadata = account.Metadata
}

func (a *rootAccountModel) domain() domain.RootAccount {
	return domain.RootAccount{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      domain.RootAccountType(a.Type),
		State:     domain.RootAccountState(a.State),
		Metadata:  a.Metadata,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (r *Repository) AddRootAccount(ctx context.Context, account domain.RootAccount) (domain.RootAccount, error) {
	db := r.DB(ctx)

	var a rootAccountModel
	a.model(account)

	now := time.Now().Unix()
	a.CreatedAt = now
	a.UpdatedAt = now

	result := db.Create(&a)
	if result.Error != nil {
		return domain.RootAccount{}, result.Error
	}

	account.CreatedAt = a.CreatedAt
	account.UpdatedAt = a.UpdatedAt

	return account, nil
}

func (r *Repository) GetRootAccountByUserID(ctx context.Context, userID uint64) (domain.RootAccount, error) {
	db := r.DB(ctx)

	var a rootAccountModel
	result := db.Unscoped().First(&a, "user_id = ?", userID)
	if result.Error != nil {
		return domain.RootAccount{}, result.Error
	}

	return a.domain(), nil
}
