package login

import (
	"bytes"
	"encoding/json"
	"github.com/sheinsviatoslav/gophermart/internal/auth"
	"github.com/sheinsviatoslav/gophermart/internal/common"
	"github.com/sheinsviatoslav/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"
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
	var credentials common.UserCredentials
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if credentials.Login == "" {
		http.Error(w, "login is required", http.StatusBadRequest)
		return
	}

	if credentials.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	user, err := h.storage.GetUserByLogin(r.Context(), credentials.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:  "userID",
		Value: user.ID,
	}

	if err := auth.WriteEncryptedCookie(w, cookie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("successful authorization")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
