package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adhamelsaady/digital-wallet/internal/ledger"
	"github.com/adhamelsaady/digital-wallet/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)
type AccountHandler struct {
	accountRepository *storage.AccountRepository
	ledgerService *ledger.Service
}

func NewAccountHandler(accountRepository *storage.AccountRepository , ledgerService *ledger.Service) *AccountHandler {
	return &AccountHandler{accountRepository: accountRepository , ledgerService: ledgerService}
}

type createAccountRequest struct {
	OwnerId uuid.UUID `json:"owner_id"`
	Currency string `json:"currency"`
	Type string `json:"type"`
}

type balanceResponse struct {
	AccountId uuid.UUID `json:"account_id"`
	Currency string `json:"currency"`
	Balance int64 `json:"balance"`
}

func (accountHandler *AccountHandler) CreateAccount (writer http.ResponseWriter , request *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid request body: malformed JSON")
		return
	}
	account , err := ledger.NewAccount(req.OwnerId, req.Currency, req.Type)
	if err != nil{
		switch {
		case errors.Is(err, ledger.ErrorInvalidCurrency), errors.Is(err, ledger.ErrorInvalidOwner):
			writeError(writer, http.StatusBadRequest, err.Error())
		default:
			writeError(writer, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	err = accountHandler.accountRepository.CreateAccount(request.Context() , account)
	if err != nil {
		if errors.Is(err, ledger.ErrorAccountDuplicate) {
			writeError(writer, http.StatusConflict, err.Error())
			return
		}
		writeError(writer, http.StatusInternalServerError, "failed to create account")
		return
	}
	writeJSON(writer, http.StatusCreated, account)
}

func (accountHandler *AccountHandler) GetBalance (writer http.ResponseWriter , request *http.Request) {
	idStr := chi.URLParam(request , "id")
	accountId , err := uuid.Parse(idStr)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid account ID")
		return
	}
	account , balance , err := accountHandler.ledgerService.GetAccountAndBalance(request.Context() , accountId)

	if err != nil {
		if errors.Is(err, ledger.ErrorAccountNotFound) {
			writeError(writer, http.StatusNotFound, err.Error())
			return
		}
		writeError(writer, http.StatusInternalServerError, "failed to get balance")
		return
	}
	response := balanceResponse {
		AccountId: accountId,
		Currency: account.Currency,
		Balance: balance,
	}
	writeJSON(writer, http.StatusOK, response)

}