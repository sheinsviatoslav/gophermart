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

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	var user common.User
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		http.Error(w, "login is required", http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	password, err := h.storage.GetUserPasswordByLogin(req.Context(), user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:  "login",
		Value: user.Login,
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
