package storage

import (
	"context"
	"encoding/json"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"time"
)

type transactionModel struct {
	ID             uint64 `gorm:"primary_key"`
	ParentID       uint64
	SrcAccountID   uint64
	DstAccountID   uint64
	SrcUserID      uint64
	DstUserID      uint64
	Amount         int64
	Currency       string
	SrcBalance     int64
	DstBalance     int64
	SrcAccountType string
	DstAccountType string
	Reversible     bool
	Metadata       []byte
	CreatedAt      int64
}

func (t *transactionModel) TableName() string {
	return "transactions"
}

func (t *transactionModel) model(transaction domain.Transaction) error {
	t.ID = transaction.ID
	t.ParentID = transaction.ParentID
	t.SrcAccountID = transaction.SrcAccountID
	t.DstAccountID = transaction.DstAccountID
	t.SrcUserID = transaction.SrcUserID
	t.DstUserID = transaction.DstUserID
	t.Amount = transaction.Amount
	t.Currency = transaction.Currency
	t.SrcBalance = transaction.SrcBalance
	t.DstBalance = transaction.DstBalance
	t.SrcAccountType = transaction.SrcAccountType
	t.DstAccountType = transaction.DstAccountType
	t.Reversible = transaction.Reversible

	meta, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return err
	}
	t.Metadata = meta

	return nil
}

func (t *transactionModel) domain() (domain.Transaction, error) {
	var metadata domain.TransactionMetadata
	if err := json.Unmarshal(t.Metadata, &metadata); err != nil {
		return domain.Transaction{}, err
	}

	return domain.Transaction{
		ID:             t.ID,
		ParentID:       t.ParentID,
		SrcAccountID:   t.SrcAccountID,
		DstAccountID:   t.DstAccountID,
		SrcUserID:      t.SrcUserID,
		DstUserID:      t.DstUserID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		SrcBalance:     t.SrcBalance,
		DstBalance:     t.DstBalance,
		SrcAccountType: t.SrcAccountType,
		DstAccountType: t.DstAccountType,
		Reversible:     t.Reversible,
		Metadata:       metadata,
		CreatedAt:      t.CreatedAt,
	}, nil
}

func (r *Repository) AddTransaction(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	db := r.DB(ctx)

	var t transactionModel
	err := t.model(transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	t.CreatedAt = time.Now().Unix()

	err = db.Create(&t).Error
	if err != nil {
		return domain.Transaction{}, wrapDBError(err, "transaction")
	}

	transaction.CreatedAt = t.CreatedAt

	return transaction, nil
}

func (r *Repository) GetTransactionByParentID(ctx context.Context, parentID uint64) (domain.Transaction, error) {
	db := r.DB(ctx)

	var t transactionModel
	result := db.Unscoped().Where("parent_id = ?", parentID).Order("created_at desc").Last(&t)
	if result.Error != nil {
		return domain.Transaction{}, result.Error
	}

	transaction, err := t.domain()
	if err != nil {
		return domain.Transaction{}, wrapDBError(err, "transaction")
	}

	return transaction, nil
}

// FixMe: Temporary for demo, delete later
func (r *Repository) GetTransactionByUserID(ctx context.Context, userID uint64) ([]domain.Transaction, error) {
	db := r.DB(ctx)

	var tm []transactionModel
	result := db.Unscoped().Where("src_user_id = ?", userID).Or("dst_user_id = ?", userID).Find(&tm)
	if result.Error != nil {
		return []domain.Transaction{}, result.Error
	}

	t := make([]domain.Transaction, len(tm))
	for i, v := range tm {
		td, err := v.domain()
		if err != nil {
			return []domain.Transaction{}, err
		}
		t[i] = td
	}

	return t, nil
}
