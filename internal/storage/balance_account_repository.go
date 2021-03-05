package storage

import (
	"context"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"time"
)

type balanceAccountModel struct {
	ID            uint64 `gorm:"primary_key"`
	RootAccountID uint64
	Type          string
	State         string
	Visible       bool
	Lockable      bool
	Balance       int64
	Currency      string
	Version       int32
	Metadata      []byte
	CreatedAt     int64
	UpdatedAt     int64
}

func (a *balanceAccountModel) TableName() string {
	return "balance_accounts"
}

func (a *balanceAccountModel) model(account domain.BalanceAccount) {
	a.ID = account.ID
	a.RootAccountID = account.RootAccountID
	a.Type = account.Type.String()
	a.State = account.State.String()
	a.Visible = account.Visible
	a.Lockable = account.Lockable
	a.Balance = account.Balance
	a.Currency = account.Currency
	a.Version = account.Version
	a.Metadata = account.Metadata
}

func (a *balanceAccountModel) domain() domain.BalanceAccount {
	return domain.BalanceAccount{
		ID:            a.ID,
		RootAccountID: a.RootAccountID,
		Type:          domain.BalanceAccountType(a.Type),
		State:         domain.BalanceAccountState(a.State),
		Visible:       a.Visible,
		Lockable:      a.Lockable,
		Balance:       a.Balance,
		Currency:      a.Currency,
		Version:       a.Version,
		Metadata:      a.Metadata,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func (r *Repository) AddBalanceAccount(ctx context.Context, account domain.BalanceAccount) (domain.BalanceAccount, error) {
	db := r.DB(ctx)

	var a balanceAccountModel
	a.model(account)

	now := time.Now().Unix()
	a.CreatedAt = now
	a.UpdatedAt = now

	err := db.Create(&a).Error
	if err != nil {
		return domain.BalanceAccount{}, wrapDBError(err, "account")
	}

	account.CreatedAt = a.CreatedAt
	account.UpdatedAt = a.UpdatedAt

	return account, nil
}

func (r *Repository) GetBalanceAccount(ctx context.Context, accountID uint64) (domain.BalanceAccount, error) {
	db := r.DB(ctx)

	var a balanceAccountModel
	err := db.Unscoped().First(&a, accountID).Error
	if err != nil {
		return domain.BalanceAccount{}, wrapDBError(err, "account")
	}

	balanceAccount := a.domain()

	return balanceAccount, nil
}

func (r *Repository) UpdateBalanceAccount(ctx context.Context, account domain.BalanceAccount) (domain.BalanceAccount, error) {
	db := r.DB(ctx)

	var a balanceAccountModel
	a.model(account)

	a.UpdatedAt = time.Now().Unix()

	err := db.Model(&a).Updates(&a).Error
	if err != nil {
		return domain.BalanceAccount{}, wrapDBError(err, "account")
	}

	account.UpdatedAt = a.UpdatedAt

	return account, nil
}

func (r *Repository) GetBalanceAccountsByRootAccountIDAndVisible(ctx context.Context, rootAccountID uint64, visible bool) ([]domain.BalanceAccount, error) {
	db := r.DB(ctx)

	var as []balanceAccountModel
	err := db.Unscoped().Where("root_account_id = ? AND visible = ?", rootAccountID, visible).Find(&as).Error
	if err != nil {
		return nil, wrapDBError(err, "account")
	}

	var accounts []domain.BalanceAccount
	for _, a := range as {
		account := a.domain()
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *Repository) GetBalanceAccountByRootAccountIDAndType(ctx context.Context, rootAccountID uint64, accountType domain.BalanceAccountType) (domain.BalanceAccount, error) {
	db := r.DB(ctx)

	var a balanceAccountModel
	err := db.Unscoped().Where("root_account_id = ? AND type = ?", rootAccountID, accountType).First(&a).Error
	if err != nil {
		return domain.BalanceAccount{}, wrapDBError(err, "account")
	}

	return a.domain(), nil
}
