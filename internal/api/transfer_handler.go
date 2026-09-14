package api

import (
	"github.com/adhamelsaady/digital-wallet/internal/ledger"
	"encoding/json"
	"net/http"
	"errors"
	"github.com/google/uuid"
)

type TransferHandler struct {
	service *ledger.Service
}

func NewTransferHandler(service *ledger.Service) *TransferHandler {
	return &TransferHandler{service: service}
}

type createTransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID uuid.UUID `json:"to_account_id"`
	Amount int64 `json:"amount"`
	Currency string `json:"currency"`
	Description string `json:"description"`
}

func (h *TransferHandler) CreateTransfer(writer http.ResponseWriter, request *http.Request) {
	var req createTransferRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid request body: malformed JSON")
		return
	}
	params := ledger.TransferParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	result , err := h.service.CreateTransfer(request.Context(), params)

	if err != nil {
		switch {
			case errors.Is(err, ledger.ErrorInvalidAmount):
			writeError(writer, http.StatusBadRequest, "amount must be greater than zero")
		case errors.Is(err, ledger.ErrorSelfTransfer):
			writeError(writer, http.StatusBadRequest, "cannot transfer to the same account")
		case errors.Is(err, ledger.ErrorAccountNotFound):
			writeError(writer, http.StatusNotFound, "one or both accounts not found")
		case errors.Is(err, ledger.ErrorCurrencyMismatch):
			writeError(writer, http.StatusConflict, "accounts must share the same currency")
		case errors.Is(err, ledger.ErrorInsufficientFunds):
			writeError(writer, http.StatusUnprocessableEntity, "insufficient funds")
		default:
			writeError(writer, http.StatusInternalServerError, "transfer failed")
		}
		return
	}
	writeJSON(writer , http.StatusCreated , result)
}
