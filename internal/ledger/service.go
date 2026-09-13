package ledger

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AccountStore interface {
	GetAccountById(ctx context.Context, accountId uuid.UUID) (*Account, error)
	CalculateBalance(ctx context.Context, accountId uuid.UUID) (int64, error)
}

type Service struct {
	accountRepository AccountStore
}

func NewService(accountRepository AccountStore) *Service {
	return &Service{accountRepository: accountRepository}
}

func (s *Service) GetAccountAndBalance(ctx context.Context, accountId uuid.UUID) (*Account, int64, error) {
	account, err := s.accountRepository.GetAccountById(ctx, accountId)
	if err != nil {
		return nil, 0, err
	}
	balance, err := s.accountRepository.CalculateBalance(ctx, accountId)
	if err != nil {
		return nil, 0, fmt.Errorf("ledger.Service: %w", err)
	}
	return account, balance, nil
}
