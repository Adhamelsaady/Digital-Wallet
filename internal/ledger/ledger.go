package ledger

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

type EntryType string
const (
	EntryTypeCredit EntryType = "CREDIT"
	EntryTypeDebit EntryType = "DEBIT"
)

var ErrorInvalidAmount = errors.New("Invalid amount")
var ErrorInvalidEntryType = errors.New("Invalid entry type")

type LedgerEntry struct {
	ID uuid.UUID `json:"id"`
	TransactionId uuid.UUID `json:"transaction_id"`
	AccountId uuid.UUID `json:"account_id"`
	Amount int64 `json:"ammount"`
	EntryType EntryType `json:"entry_type"`
	CreatedAt time.Time `json:"time"` 
}

func CalculateBalance(entries []LedgerEntry) int64 {
	var balance int64 = 0
	for _, entry := range entries {
		if entry.EntryType == EntryTypeCredit {
			balance += entry.Amount
		} else {
			balance -= entry.Amount
		}
	}
	return balance
}