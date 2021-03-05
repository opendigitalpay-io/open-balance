package idem

import (
	"context"
	"errors"
	"github.com/opendigitalpay-io/open-balance/internal/common/uid"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"github.com/opendigitalpay-io/open-balance/internal/storage"
)

type Service interface {
	Start(context.Context, string) error
	End(context.Context, string, interface{}) error
	IdemExec(context.Context, string, func() (interface{}, error)) (interface{}, error)
}

type Repository interface {
	AddIdemEntry(context.Context, uint64, string) error
	UpdateIdemEntry(context.Context, string, interface{}) error
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

func (s *service) Start(ctx context.Context, idemID string) error {
	id, idGenErr := s.uidGenerator.NextID()
	if idGenErr != nil {
		return idGenErr
	}

	err := s.repo.AddIdemEntry(ctx, id, idemID)
	if err != nil {
		var dee storage.DuplicatedEntryError
		if errors.As(err, &dee) {
			return domain.IdemError{
				What: "Idempotency Key In Use",
			}
		}
		return err
	}
	return nil
}

func (s *service) End(ctx context.Context, idemID string, response interface{}) error {
	err := s.repo.UpdateIdemEntry(ctx, idemID, response)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) IdemExec(ctx context.Context, idemID string, bizF func() (interface{}, error)) (interface{}, error) {
	errStartF := s.Start(ctx, idemID)
	if errStartF != nil {
		return nil, errStartF
	}

	response, errBizF := bizF()
	if errBizF != nil {
		return nil, errBizF
	}

	errEndF := s.End(ctx, idemID, response)
	if errEndF != nil {
		return nil, errEndF
	}
	return response, nil
}
