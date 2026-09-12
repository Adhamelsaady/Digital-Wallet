package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/google/uuid"
	"github.com/adhamelsaady/digital-wallet/internal/ledger"
	"github.com/adhamelsaady/digital-wallet/internal/storage"
)
type AccountHandler struct {
	accountRepository *storage.AccountRepository
}

func NewAccountHandler(repository *storage.AccountRepository) *AccountHandler {
	return &AccountHandler{accountRepository: repository}
}

type createAccountRequest struct {
	OwnerId uuid.UUID `json:"owner_id"`
	Currency string `json:"currency"`
	Type string `json:"type"`
}

type balanceResponse struct {
	AccountId uuid.UUID `json:"account_id"`
	Currency string `json:"currency"`
	Balance string `json:"balance"`
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

func (accountHandler *AccountHandler) GetAccountBalance(){
	// to-do
}