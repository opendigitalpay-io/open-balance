package pay

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
	TryExternalPay(context.Context, string, api.TryPayRequest) (uint64, error)
	TryPay(context.Context, string, api.TryPayRequest) (uint64, error)
	CommitPay(context.Context, api.CommitPayRequest) error
	CancelPay(context.Context, api.CancelPayRequest) error
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

func (s *service) TryExternalPay(ctx context.Context, idemKey string, req api.TryPayRequest) (uint64, error) {
	parentID, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return s.idemService.IdemExec(ctxWithTxn, idemKey, func() (interface{}, error) {
			return s.tryExternalPay(ctxWithTxn, req)
		})
	})

	if err != nil {
		return 0, err
	}

	return parentID.(uint64), nil
}

func (s *service) TryPay(ctx context.Context, idemKey string, req api.TryPayRequest) (uint64, error) {
	parentID, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return s.idemService.IdemExec(ctxWithTxn, idemKey, func() (interface{}, error) {
			return s.tryPay(ctxWithTxn, req)
		})
	})

	if err != nil {
		return 0, err
	}

	return parentID.(uint64), nil
}

func (s *service) CommitPay(ctx context.Context, req api.CommitPayRequest) error {
	_, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return nil, s.commitPay(ctxWithTxn, req)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CancelPay(ctx context.Context, req api.CancelPayRequest) error {
	_, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		return nil, s.cancelPay(ctxWithTxn, req)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) tryExternalPay(ctx context.Context, req api.TryPayRequest) (uint64, error) {
	parentID, err := s.uidGenerator.NextID()

	if err != nil {
		return 0, err
	}

	// Here, we are using userId 1 as system user.
	systemRootAccount, err := s.repo.GetRootAccountByUserID(ctx, 1)
	if err != nil {
		return 0, err
	}
	systemSourceAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, systemRootAccount.ID, domain.SOURCE)
	if err != nil {
		return 0, nil
	}

	customerRootAccount, err := s.repo.GetRootAccountByUserID(ctx, req.UserID)
	if err != nil {
		return 0, err
	}
	customerPaymentAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, customerRootAccount.ID, domain.PAYMENT)
	if err != nil {
		return 0, err
	}

	merchantRootAccount, err := s.repo.GetRootAccountByUserID(ctx, req.BusinessID)
	if err != nil {
		return 0, err
	}
	merchantPayableAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, merchantRootAccount.ID, domain.PAYABLE)
	if err != nil {
		return 0, err
	}

	transactionItem := domain.TransactionItem{
		SrcAccount: systemSourceAccount,
		SrcUserID:  systemRootAccount.UserID,
		DstAccount: customerPaymentAccount,
		DstUserID:  customerRootAccount.UserID,
		Amount:     req.Amount,
	}

	metadata := domain.TransactionMetadata{
		Action: domain.TRY,
		NextTransaction: domain.NextTransaction{
			SrcAccountID: customerPaymentAccount.ID,
			SrcUserID:    customerRootAccount.UserID,
			DstAccountID: merchantPayableAccount.ID,
			DstUserID:    merchantRootAccount.UserID,
			Amount:       req.Amount,
		},
	}

	err = s.transactionService.Transfer(ctx, parentID, transactionItem, metadata)
	if err != nil {
		return 0, err
	}

	return parentID, nil
}

func (s *service) tryPay(ctx context.Context, req api.TryPayRequest) (uint64, error) {

	parentID, err := s.uidGenerator.NextID()

	if err != nil {
		return 0, err
	}

	customerRootAccount, err := s.repo.GetRootAccountByUserID(ctx, req.UserID)
	if err != nil {
		return 0, err
	}
	customerChequeAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, customerRootAccount.ID, domain.CHEQUE)
	if err != nil {
		return 0, nil
	}
	customerPaymentAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, customerRootAccount.ID, domain.PAYMENT)
	if err != nil {
		return 0, err
	}

	merchantRootAccount, err := s.repo.GetRootAccountByUserID(ctx, req.BusinessID)
	if err != nil {
		return 0, err
	}
	merchantPayableAccount, err := s.repo.GetBalanceAccountByRootAccountIDAndType(ctx, merchantRootAccount.ID, domain.PAYABLE)
	if err != nil {
		return 0, err
	}

	transactionItem := domain.TransactionItem{
		SrcAccount: customerChequeAccount,
		SrcUserID:  customerRootAccount.UserID,
		DstAccount: customerPaymentAccount,
		DstUserID:  customerRootAccount.UserID,
		Amount:     req.Amount,
	}

	metadata := domain.TransactionMetadata{
		Action: domain.TRY,
		NextTransaction: domain.NextTransaction{
			SrcAccountID: customerPaymentAccount.ID,
			SrcUserID:    customerRootAccount.UserID,
			DstAccountID: merchantPayableAccount.ID,
			DstUserID:    merchantRootAccount.UserID,
			Amount:       req.Amount,
		},
	}

	err = s.transactionService.Transfer(ctx, parentID, transactionItem, metadata)
	if err != nil {
		return 0, err
	}

	return parentID, nil
}

func (s *service) commitPay(ctx context.Context, req api.CommitPayRequest) error {
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

	meta := domain.TransactionMetadata{
		Action: domain.COMMIT,
	}
	err = s.transactionService.Transfer(ctx, req.ParentID, transactionItem, meta)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) cancelPay(ctx context.Context, req api.CancelPayRequest) error {
	prevTransaction, err := s.transactionService.GetTransactionByParentID(ctx, req.ParentID)
	if err != nil {
		return err
	}

	if !prevTransaction.Reversible || prevTransaction.Metadata.Action != domain.TRY {
		return domain.TransactionError{
			What: fmt.Sprintf("%d cannot be cancelled", req.ParentID),
		}
	}

	srcAccount, err := s.repo.GetBalanceAccount(ctx, prevTransaction.SrcAccountID)
	if err != nil {
		return err
	}

	dstAccount, err := s.repo.GetBalanceAccount(ctx, prevTransaction.DstAccountID)

	transactionItem := domain.TransactionItem{
		SrcAccount: dstAccount,
		SrcUserID:  prevTransaction.DstUserID,
		DstAccount: srcAccount,
		DstUserID:  prevTransaction.SrcUserID,
		Amount:     prevTransaction.Amount,
	}

	meta := domain.TransactionMetadata{
		Action: domain.CANCEL,
	}

	err = s.transactionService.Transfer(ctx, req.ParentID, transactionItem, meta)
	if err != nil {
		return err
	}

	return nil
}
