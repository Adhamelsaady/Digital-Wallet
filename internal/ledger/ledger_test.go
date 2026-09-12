package ledger_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/adhamelsaady/digital-wallet/internal/ledger"
)

func TestCalculateBalance(t *testing.T) {
	accountID := uuid.New()
	transactionID := uuid.New()
	now := time.Now().UTC()
	tests := []struct {
		name string
		entries []ledger.LedgerEntry
		expectedBalance int64
	}{
		{
			name:"zero entries returns zero balance",
			entries:[]ledger.LedgerEntry{},
			expectedBalance: 0,
		},
		{
			name:"single credit entry adds to balance",
			entries:[]ledger.LedgerEntry{
				{
					ID:uuid.New(),
					TransactionId:transactionID,
					AccountId:accountID,
					Amount:10000,				
					EntryType:ledger.EntryTypeCredit,
					CreatedAt:now,
				},
			},
			expectedBalance: 10000,
		},
		{
			name: "single debit entry produces negative balance",
			entries: []ledger.LedgerEntry{
				{
					ID:uuid.New(),
					TransactionId:transactionID,
					AccountId:accountID,
					Amount:2500, // $25.00
					EntryType:ledger.EntryTypeDebit,
					CreatedAt:now,
				},
			},
			expectedBalance: -2500,
		},
		{
			name: "mixed debit and credit entries sum correctly",
			entries: []ledger.LedgerEntry{
				{
					ID:uuid.New(),
					TransactionId:transactionID,
					AccountId:accountID,
					Amount:10000,
					EntryType:ledger.EntryTypeCredit,
					CreatedAt:now,
				},
				{
					ID:uuid.New(),
					TransactionId:transactionID,
					AccountId:accountID,
					Amount:3500,
					EntryType:ledger.EntryTypeDebit,
					CreatedAt:now,
				},
				{
					ID:uuid.New(),
					TransactionId:transactionID,
					AccountId:accountID,
					Amount:500,
					EntryType:ledger.EntryTypeCredit,
					CreatedAt:now,
				},
			},
			expectedBalance: 7000, 
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := ledger.CalculateBalance(tc.entries)
			assert.Equal(t, tc.expectedBalance, actual)
		})
	}
}
