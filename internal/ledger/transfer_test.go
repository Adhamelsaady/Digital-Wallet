package ledger

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAccountStore struct{}

func (m *mockAccountStore) GetAccountById(_ context.Context, _ uuid.UUID) (*Account, error) {
	return &Account{}, nil
}
func (m *mockAccountStore) CalculateBalance(_ context.Context, _ uuid.UUID) (int64, error) {
	return 0, nil
}

type mockTransferStore struct {
	called bool
}

func (m *mockTransferStore) ExecuteTransfer(_ context.Context, params TransferParams) (*TransferResponse, error) {
	m.called = true
	return &TransferResponse{
		TransactionId: uuid.New(),
		FromAccountId: params.FromAccountId,
		ToAccountId: params.ToAccountId,
		Amount: params.Amount,
		Currency: "USD",
		Status: "COMPLETED",
	}, nil
}

func TestCreateTransfer_Validation(t *testing.T) {
	tests := []struct {
		name string
		params TransferParams
		wantErr error
		wantCall bool
	}{
		{
			name: "valid transfer passes domain checks and calls store",
			params: TransferParams{
				FromAccountId: uuid.New(),
				ToAccountId: uuid.New(),
				Amount: 1000,
			},
			wantErr: nil,
			wantCall: true,
		},
		{
			name: "zero amount is rejected before reaching the store",
			params: TransferParams{
				FromAccountId: uuid.New(),
				ToAccountId: uuid.New(),
				Amount: 0,
			},
			wantErr: ErrorInvalidAmount,
			wantCall: false,
		},
		{
			name: "negative amount is rejected before reaching the store",
			params: TransferParams{
				FromAccountId: uuid.New(),
				ToAccountId: uuid.New(),
				Amount: -500,
			},
			wantErr: ErrorInvalidAmount,
			wantCall: false,
		},
		{
			name: "self-transfer is rejected before reaching the store",
			params: func() TransferParams {
				id := uuid.New()
				return TransferParams{
					FromAccountId: id,
					ToAccountId: id,
					Amount: 1000,
				}
			}(),
			wantErr: ErrorSelfTransfer,
			wantCall: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTransfer := &mockTransferStore{}
			svc := NewService(&mockAccountStore{}, mockTransfer)
			result, err := svc.CreateTransfer(context.Background(), tc.params)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			assert.Equal(t, tc.wantCall, mockTransfer.called)
		})
	}
}
