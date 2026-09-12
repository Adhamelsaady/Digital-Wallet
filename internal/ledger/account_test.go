package ledger_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adhamelsaady/digital-wallet/internal/ledger"
)

func TestNewAccount(t *testing.T) {
	validOwnerId := uuid.New()
	tests := []struct {
		description string
		ownerID uuid.UUID
		currency string
		accountType string
		intendedError error
	}{
		{
			description: "valid USD account",
			ownerID: validOwnerId,
			currency: "USD",	
			accountType: "AVAILABLE",
			intendedError: nil,
		},
		{
			description: "valid EGP account lowercase gets normalized",
			ownerID: validOwnerId,
			currency: "egp",
			accountType: "",
			intendedError: nil,
		},
		{
			description: "invalid currency length",
			ownerID: validOwnerId,
			currency: "US",
			accountType: "AVAILABLE",
			intendedError: ledger.ErrorInvalidCurrency,
		},
		{
			description: "unsupported currency",
			ownerID: validOwnerId,
			currency: "XYZ",
			accountType: "AVAILABLE",
			intendedError: ledger.ErrorInvalidCurrency,
		},
		{
			description: "Zero-value ownerID",
			ownerID: uuid.Nil,
			currency: "USD",
			accountType: "AVAILABLE",
			intendedError: ledger.ErrorInvalidOwner,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			account, err := ledger.NewAccount(testCase.ownerID, testCase.currency, testCase.accountType)

			if testCase.intendedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.intendedError)
				assert.Nil(t, account)
			} else {
				require.NoError(t, err)
				require.NotNil(t, account)
				assert.Equal(t, testCase.ownerID, account.OwnerID)
				assert.NotEmpty(t, account.Currency)
				assert.NotEmpty(t, account.Type)
			}
		})
	}
}
