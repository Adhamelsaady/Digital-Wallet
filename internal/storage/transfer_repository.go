package storage

import (
	"context"
	"fmt"
	"errors"
	"time"
	"github.com/google/uuid"

	"github.com/adhamelsaady/digital-wallet/internal/ledger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	referenceTransfer = "Transfer"
	statusCompleted = "Completed"
	entryDebit = "Debit"
	entryCredit = "Credit"
)

type TransferRepository struct {
	pool *pgxpool.Pool
}


func NewTransferRepository(pool *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{pool: pool}
}

func (r *TransferRepository) ExecuteTransfer(ctx context.Context, transferParams ledger.TransferParams) (*ledger.TransferResponse, error) {
	// begin transaction
	transaction , err := r.pool.Begin(ctx)
	if err != nil {
		return nil , fmt.Errorf("begin transaction: %w", err)
	}
	// defer rollback
	defer transaction.Rollback(ctx)
	from , err := getAccountById(ctx , transaction , transferParams.FromAccountId)
	if err != nil {
		return nil, fmt.Errorf("get sender account: %w", err)
	}
	to , err := getAccountById(ctx , transaction , transferParams.ToAccountId)
	if err != nil {
		return nil, fmt.Errorf("get receiver account: %w", err)
	}
	err = validateTransfer(from , to, transferParams.Amount)
	if err != nil {
		return nil, fmt.Errorf("validate transfer: %w", err)
	}
	balance , err := calculateBalanceTx(ctx, transaction, from.ID)
	if err != nil {
		return nil, fmt.Errorf("calculate sender balance: %w", err)
	}
	if balance < transferParams.Amount {
		return nil, ledger.ErrorInsufficientFunds
	}
	transactionId := uuid.New()
	curTime := time.Now().UTC()
	if err := createTransaction(ctx , transaction , transactionId , transferParams.Description , curTime); err != nil {
		return nil, err
	}

	if err := createLedgerEntry(ctx , transaction , transactionId, transferParams.FromAccountId , transferParams.Amount , entryDebit , curTime); err != nil {
		return nil, err
	}
	if err := createLedgerEntry(ctx , transaction , transactionId, transferParams.ToAccountId , transferParams.Amount , entryCredit , curTime); err != nil {
		return nil, err
	}
	if err := transaction.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return &ledger.TransferResponse{
		TransactionId: transactionId,
		FromAccountId: from.ID,
		ToAccountId: to.ID,
		Amount: transferParams.Amount,
		Currency: from.Currency,
		Description: transferParams.Description,
		Status: statusCompleted,
		CreatedAt: curTime,
	},nil

}
// func createLedgerEntry(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID, accountID uuid.UUID, amount int64, entryType string, createdAt time.Time) error {
func getAccountById (ctx context.Context , transaction pgx.Tx , id uuid.UUID) (*ledger.Account , error) {
	query := `SELECT id, owner_id, currency, type, created_at, updated_at FROM accounts WHERE id = $1 FOR UPDATE`
	var account ledger.Account
	err := transaction.QueryRow(ctx , query , id).Scan(
		&account.ID,
		&account.OwnerID,
		&account.Currency,
		&account.Type,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ledger.ErrorAccountNotFound
		}

		return nil, fmt.Errorf(
			"get account %s: %w",
			id,
			err,
		) 
	}
	return &account, nil
}
func validateTransfer(from *ledger.Account,to *ledger.Account,amount int64) error {
	if from.ID == to.ID {
		return ledger.ErrorSelfTransfer
	}
	if amount <= 0 {
		return ledger.ErrorInvalidAmount
	}
	if from.Currency != to.Currency {
		return ledger.ErrorCurrencyMismatch
	}
	return nil
}

func calculateBalanceTx(ctx context.Context, tx pgx.Tx, accountID uuid.UUID) (int64, error) {
	const query = `
		SELECT COALESCE(
			SUM(
				CASE
					WHEN entry_type = 'CREDIT' THEN amount
					WHEN entry_type = 'DEBIT'  THEN -amount
					ELSE 0
				END
			),
			0
		)
		FROM ledger_entries	WHERE account_id = $1
	`
	var balance int64
	err := tx.QueryRow(ctx,query,accountID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("calculate balance for %s: %w",accountID,err)
	}
	return balance, nil
}


func createTransaction(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID, description string, createdAt time.Time) error {
	const query = `INSERT INTO transactions (id,reference_type,description,status,	created_at)	VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, transactionID, referenceTransfer, description, statusCompleted, createdAt)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func createLedgerEntry(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID, accountID uuid.UUID, amount int64, entryType string, createdAt time.Time) error {
	const query = `INSERT INTO ledger_entries(id, transaction_id, account_id, amount, entry_type, created_at) 
				   VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := tx.Exec(ctx, query,uuid.New(), transactionID, accountID, amount, entryType, createdAt)
	if err != nil {
		return fmt.Errorf("create %s ledger entry: %w",	entryType, err)
	}
	return nil
}