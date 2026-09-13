package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/adhamelsaady/digital-wallet/internal/ledger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) CreateAccount(ctx context.Context,
	account *ledger.Account) error {
	query := `INSERT INTO accounts (id, owner_id, currency, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query,
		account.ID, account.OwnerID,
		account.Currency, account.Type,
		account.CreatedAt, account.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ledger.ErrorAccountDuplicate
		}
		return fmt.Errorf("storage: failed to insert account: %w", err)
	}
	return err
}

func (r *AccountRepository) GetAccountById(ctx context.Context, id uuid.UUID) (*ledger.Account, error) {
	query := `SELECT id, owner_id, currency, type, created_at, updated_at FROM accounts WHERE id = $1`
	var result ledger.Account
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&result.ID, &result.OwnerID, &result.Currency, &result.Type,
		&result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ledger.ErrorAccountNotFound
		}
		return nil, fmt.Errorf("storage: failed to get account %s: %w", id, err)
	}
	return &result, nil
}

func (r *AccountRepository) CalculateBalance(ctx context.Context, accountId uuid.UUID) (int64, error) {

	query := `SELECT COALESCE(
			SUM(
				CASE 
					WHEN entry_type = 'CREDIT' THEN amount 
					WHEN entry_type = 'DEBIT'  THEN -amount 
					ELSE 0 
				END
			), 0 )FROM ledger_entries WHERE account_id = $1`

	var balance int64
	err := r.pool.QueryRow(ctx, query, accountId).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("storage: failed to calculate balance for account %s: %w", accountId, err)
	}
	return balance, nil
}
