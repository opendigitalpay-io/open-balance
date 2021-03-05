package topup

import (
	"context"
	"fmt"
	"github.com/opendigitalpay-io/open-balance/internal/common/uid"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"github.com/opendigitalpay-io/open-balance/internal/idem"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
	"github.com/opendigitalpay-io/open-balance/internal/transaction"
)

type Service interface {
	TryTopUp(context.Context, string, api.TryTopUpRequest) (uint64, error)
	CommitTopUp(context.Context, api.CommitTopUpRequest) error
	CancelTopUp(context.Context, api.CancelTopUpRequest) error
}

type Repository interface {
	GetRootAccountByUserID(context.Context, uint64) (domain.RootAccount, error)

	GetBalanceAccount(context.Context, uint64) (domain.BalanceAccount, error)
	GetBalanceAccountByRootAccountIDAndType(context.Context, uint64, domain.BalanceAccountType) (domain.BalanceAccount, error)

	TxnExec(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
}

type service struct {
	repo               Repository
	uidGenerator       uid.Generator
	transactionService transaction.Service
	idemService        idem.Service
}

func NewService(repo Repository, uidGenerator uid.Generator, transactionService transaction.Service, idemService idem.Service) Service {
	return &service{
		repo:               repo,
		uidGenerator:       uidGenerator,
		transactionService: transactionService,
		idemService:        idemService,
	}
}

func (s *service) TryTopUp(ctx context.Context, idemID string, req api.TryTopUpRequest) (uint64, error) {
	parentID, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return s.idemService.IdemExec(ctxWithTxn, idemID, func() (interface{}, error) {
			return s.tryTopUp(ctxWithTxn, req)
		})
	})

	if err != nil {
		return 0, err
	}

	return parentID.(uint64), nil
}

func (s *service) CommitTopUp(ctx context.Context, req api.CommitTopUpRequest) error {
	_, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return nil, s.commitTopUp(ctxWithTxn, req)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CancelTopUp(ctx context.Context, req api.CancelTopUpRequest) error {
	_, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return nil, s.cancelTopUp(ctxWithTxn, req)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) tryTopUp(ctx context.Context, req api.TryTopUpRequest) (uint64, error) {
	parentID, err := s.uidGenerator.NextID()
	if err != nil {
		return 0, err
	}

	systemRootAccount, err := s.repo.GetRootAccountByUserID(ctx, 1)
	if err != nil {
		return 0, err
	}
	systemSourceAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, systemRootAccount.ID, domain.SOURCE)
	if err != nil {
		return 0, err
	}

	customerRootAccount, err := s.repo.GetRootAccountByUserID(ctx, req.UserID)
	if err != nil {
		return 0, err
	}
	customerIncomingAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, customerRootAccount.ID, domain.INCOMING)
	if err != nil {
		return 0, err
	}
	customerChequeAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, customerRootAccount.ID, domain.CHEQUE)
	if err != nil {
		return 0, err
	}

	transactionItem := domain.TransactionItem{
		SrcAccount: systemSourceAccount,
		SrcUserID:  systemRootAccount.UserID,
		DstAccount: customerIncomingAccount,
		DstUserID:  customerRootAccount.UserID,
		Amount:     req.Amount,
	}

	metadata := domain.TransactionMetadata{
		Action: domain.TRY,
		NextTransaction: domain.NextTransaction{
			SrcAccountID: customerIncomingAccount.ID,
			SrcUserID:    customerRootAccount.UserID,
			DstAccountID: customerChequeAccount.ID,
			DstUserID:    customerRootAccount.ID,
			Amount:       req.Amount,
		},
	}

	err = s.transactionService.Transfer(ctx, parentID, transactionItem, metadata)
	if err != nil {
		return 0, err
	}

	return parentID, nil
}

func (s *service) commitTopUp(ctx context.Context, req api.CommitTopUpRequest) error {
	prevTransaction, err := s.transactionService.GetTransactionByParentID(ctx, req.ParentID)
	if err != nil {
		return err
	}

	if prevTransaction.Metadata.Action != domain.TRY {
		return domain.TransactionError{
			What: fmt.Sprintf("%d cannot be committed", req.ParentID),
		}
	}

	nextTransaction := prevTransaction.Metadata.NextTransaction

	srcAccount, err := s.repo.GetBalanceAccount(ctx, nextTransaction.SrcAccountID)
	if err != nil {
		return err
	}

	dstAccount, err := s.repo.GetBalanceAccount(ctx, nextTransaction.DstAccountID)
	if err != nil {
		return err
	}

	transactionItem := domain.TransactionItem{
		SrcAccount: srcAccount,
		SrcUserID:  nextTransaction.SrcUserID,
		DstAccount: dstAccount,
		DstUserID:  nextTransaction.DstUserID,
		Amount:     nextTransaction.Amount,
	}

	metadata := domain.TransactionMetadata{
		Action: domain.COMMIT,
	}

	err = s.transactionService.Transfer(ctx, req.ParentID, transactionItem, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) cancelTopUp(ctx context.Context, req api.CancelTopUpRequest) error {
	prevTransaction, err := s.transactionService.GetTransactionByParentID(ctx, req.ParentID)
	if err != nil {
		return err
	}

	if prevTransaction.Metadata.Action != domain.TRY || !prevTransaction.Reversible {
		return domain.TransactionError{
			What: fmt.Sprintf("%d cannot be cancelled", req.ParentID),
		}
	}

	srcAccount, err := s.repo.GetBalanceAccount(ctx, prevTransaction.SrcAccountID)
	if err != nil {
		return err
	}

	dstAccount, err := s.repo.GetBalanceAccount(ctx, prevTransaction.DstAccountID)
	if err != nil {
		return err
	}

	transactionItem := domain.TransactionItem{
		SrcAccount: dstAccount,
		SrcUserID:  prevTransaction.DstUserID,
		DstAccount: srcAccount,
		DstUserID:  prevTransaction.SrcAccountID,
		Amount:     prevTransaction.Amount,
	}

	metadata := domain.TransactionMetadata{
		Action: domain.CANCEL,
	}

	err = s.transactionService.Transfer(ctx, req.ParentID, transactionItem, metadata)
	if err != nil {
		return err
	}

	return nil
}
