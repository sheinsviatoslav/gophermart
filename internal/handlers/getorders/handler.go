package getorders

import (
	"encoding/json"
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
	orders, err := h.storage.GetUserOrders(r.Context(), utils.GetUserID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		http.Error(w, "no orders from user", http.StatusNoContent)
		return
	}

	response, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
