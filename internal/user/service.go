package user

import (
	"context"
	"github.com/opendigitalpay-io/open-balance/internal/common/uid"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
)

type Service interface {
	AddUser(context.Context, api.AddUserRequest) (User, error)
	GetUser(context.Context, uint64) (User, error)
	GetVisibleBalanceAccounts(context.Context, uint64) ([]domain.BalanceAccount, error)
	UpdateUser(context.Context, uint64, api.UpdateUserRequest) (User, error)
}

type Repository interface {
	AddUser(context.Context, User) (User, error)
	GetUser(context.Context, uint64) (User, error)
	UpdateUser(context.Context, User) (User, error)

	AddRootAccount(context.Context, domain.RootAccount) (domain.RootAccount, error)
	GetRootAccountByUserID(context.Context, uint64) (domain.RootAccount, error)

	AddBalanceAccount(context.Context, domain.BalanceAccount) (domain.BalanceAccount, error)
	GetBalanceAccountsByRootAccountIDAndVisible(context.Context, uint64, bool) ([]domain.BalanceAccount, error)

	TxnExec(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
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

func (s *service) AddUser(ctx context.Context, req api.AddUserRequest) (User, error) {
	u, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		userID, err := s.uidGenerator.NextID()
		if err != nil {
			return User{}, err
		}

		user := User{
			ID:         userID,
			Email:      req.Email,
			Phone:      req.Phone,
			ExternalID: req.UserName,
		}

		user, err = s.repo.AddUser(ctxWithTxn, user)
		if err != nil {
			return User{}, err
		}

		rootAccountID, err := s.uidGenerator.NextID()
		if err != nil {
			return User{}, err
		}
		rootAccount := domain.NewRootAccount(rootAccountID, userID, domain.PERSONAL)
		_, err = s.repo.AddRootAccount(ctxWithTxn, rootAccount)
		if err != nil {
			return User{}, err
		}

		// TODO: currency
		chequeAccountID, err := s.uidGenerator.NextID()
		if err != nil {
			return User{}, err
		}
		chequeAccount := domain.NewBalanceAccount(chequeAccountID, rootAccountID, domain.CHEQUE, true, true, "CAD")
		_, err = s.repo.AddBalanceAccount(ctxWithTxn, chequeAccount)
		if err != nil {
			return User{}, err
		}
		incomingAccountID, err := s.uidGenerator.NextID()
		if err != nil {
			return User{}, err
		}
		incomingAccount := domain.NewBalanceAccount(incomingAccountID, rootAccountID, domain.INCOMING, false, true, "CAD")
		_, err = s.repo.AddBalanceAccount(ctxWithTxn, incomingAccount)
		if err != nil {
			return User{}, err
		}
		paymentAccountID, err := s.uidGenerator.NextID()
		if err != nil {
			return User{}, err
		}
		paymentAccount := domain.NewBalanceAccount(paymentAccountID, rootAccountID, domain.PAYMENT, true, true, "CAD")
		_, err = s.repo.AddBalanceAccount(ctxWithTxn, paymentAccount)
		if err != nil {
			return User{}, err
		}

		return user, nil
	})

	if err != nil {
		return User{}, err
	}

	return u.(User), nil
}

func (s *service) GetUser(ctx context.Context, userID uint64) (User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *service) GetVisibleBalanceAccounts(ctx context.Context, userID uint64) ([]domain.BalanceAccount, error) {
	rootAccount, err := s.repo.GetRootAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	visibleBalanceAccounts, err := s.repo.GetBalanceAccountsByRootAccountIDAndVisible(ctx, rootAccount.ID, true)
	if err != nil {
		return nil, err
	}

	return visibleBalanceAccounts, nil
}

func (s *service) UpdateUser(ctx context.Context, userID uint64, req api.UpdateUserRequest) (User, error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return User{}, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.UserName != "" {
		user.ExternalID = req.UserName
	}

	user, err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
