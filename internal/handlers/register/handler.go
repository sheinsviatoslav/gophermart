package register

import (
	"bytes"
	"encoding/json"
	"github.com/sheinsviatoslav/gophermart/internal/auth"
	"github.com/sheinsviatoslav/gophermart/internal/common"
	"github.com/sheinsviatoslav/gophermart/internal/storage"
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

	if len(user.Login) < 6 {
		http.Error(w, "login is too short", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 8 {
		http.Error(w, "password is too short", http.StatusBadRequest)
		return
	}

	loginExists, err := h.storage.CheckLoginExists(req.Context(), user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loginExists {
		http.Error(w, "user with this login already exists", http.StatusConflict)
		return
	}

	if err = h.storage.CreateUser(req.Context(), common.User{
		Login:    user.Login,
		Password: user.Password,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if _, err := w.Write([]byte("user successfully registered")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
