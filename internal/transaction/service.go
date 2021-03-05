package transaction

import (
	"context"
	"github.com/opendigitalpay-io/open-balance/internal/common/uid"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
)

type Service interface {
	Transfer(context.Context, uint64, domain.TransactionItem, domain.TransactionMetadata) error
	GetTransactionByParentID(context.Context, uint64) (domain.Transaction, error)

	// FixMe: Temporary for demo, delete later
	GetTransactionByUserID(context.Context, uint64) ([]domain.Transaction, error)
}

type Repository interface {
	UpdateBalanceAccount(context.Context, domain.BalanceAccount) (domain.BalanceAccount, error)

	AddTransaction(context.Context, domain.Transaction) (domain.Transaction, error)
	GetTransactionByParentID(context.Context, uint64) (domain.Transaction, error)

	// FixMe: Temporary for demo, delete later
	GetTransactionByUserID(context.Context, uint64) ([]domain.Transaction, error)
}

type service struct {
	repo         Repository
	uidGenerator uid.Generator
}

func NewService(repo Repository, uidGenerator uid.Generator) Service {
	return &service{
		repo:         repo,
		uidGenerator: uidGenerator,
	}
}

// FixMe: Temporary for demo, delete later
func (s *service) GetTransactionByUserID(ctx context.Context, userID uint64) ([]domain.Transaction, error) {
	transactions, err := s.repo.GetTransactionByUserID(ctx, userID)
	if err != nil {
		return []domain.Transaction{}, nil
	}

	return transactions, nil
}

func (s *service) GetTransactionByParentID(ctx context.Context, parentID uint64) (domain.Transaction, error) {
	transaction, err := s.repo.GetTransactionByParentID(ctx, parentID)
	if err != nil {
		return domain.Transaction{}, err
	}

	return transaction, nil
}

func (s *service) Transfer(ctx context.Context, parentID uint64, transactionItem domain.TransactionItem, meta domain.TransactionMetadata) error {

	fromAccount := transactionItem.SrcAccount
	toAccount := transactionItem.DstAccount
	amount := transactionItem.Amount
	if fromAccount.Currency != toAccount.Currency {
		return domain.TransactionError{
			What: "transaction currency doesn't match",
		}
	}

	if fromAccount.Balance < amount {
		return domain.TransactionError{
			What: "No sufficient funds",
		}
	}

	fromAccount.Debit(amount)
	_, err := s.repo.UpdateBalanceAccount(ctx, fromAccount)
	if err != nil {
		return err
	}

	toAccount.Credit(amount)
	_, err = s.repo.UpdateBalanceAccount(ctx, toAccount)
	if err != nil {
		return err
	}

	transactionID, err := s.uidGenerator.NextID()
	if err != nil {
		return err
	}

	transaction := domain.Transaction{
		ID:             transactionID,
		ParentID:       parentID,
		SrcAccountID:   fromAccount.ID,
		DstAccountID:   toAccount.ID,
		SrcUserID:      transactionItem.SrcUserID,
		DstUserID:      transactionItem.DstUserID,
		Amount:         amount,
		Currency:       toAccount.Currency,
		SrcBalance:     fromAccount.Balance,
		DstBalance:     toAccount.Balance,
		SrcAccountType: fromAccount.Type.String(),
		DstAccountType: toAccount.Type.String(),
		Reversible:     true,
		Metadata:       meta,
	}

	_, err = s.repo.AddTransaction(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}
