package addwithdrawal

import (
	"bytes"
	"encoding/json"
	"github.com/sheinsviatoslav/gophermart/internal/common"
	"github.com/sheinsviatoslav/gophermart/internal/storage"
	"github.com/sheinsviatoslav/gophermart/internal/utils"
	"net/http"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var withdrawal common.Withdrawal
	var buf bytes.Buffer
	userID := utils.GetUserID(r)

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &withdrawal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := h.storage.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if balance.Current < withdrawal.Sum {
		http.Error(w, "insufficient balance", http.StatusPaymentRequired)
		return
	}

	if !utils.IsOrderNumberValid(withdrawal.Order) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	if err = h.storage.AddWithdrawal(r.Context(), withdrawal.Sum, withdrawal.Order, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("success")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
