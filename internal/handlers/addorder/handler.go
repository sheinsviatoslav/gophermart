package addorder

import (
	"context"
	"fmt"
	"github.com/sheinsviatoslav/gophermart/internal/storage"
	"github.com/sheinsviatoslav/gophermart/internal/utils"
	"io"
	"net/http"
	"time"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber := string(body)

	if !utils.IsOrderNumberValid(orderNumber) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	order, err := h.storage.GetOrderByNumber(r.Context(), orderNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := utils.GetUserID(r)

	if order.UserID == userID {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("success")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if order.UserID != "" {
		http.Error(w, "order number is already created by another user", http.StatusConflict)
		return
	}

	if err = h.storage.AddOrder(r.Context(), orderNumber, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if err = h.storage.UpdateOrderFromAccrual(context.Background(), done, orderNumber, userID); err != nil {
					fmt.Printf("error updating order from accrual: %v\n", err)
					return
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte("order successfully added")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
