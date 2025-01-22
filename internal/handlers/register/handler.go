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

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var user common.User
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
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

	loginExists, err := h.storage.CheckLoginExists(r.Context(), user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loginExists {
		http.Error(w, "user with this login already exists", http.StatusConflict)
		return
	}

	newUser, err := h.storage.CreateUser(r.Context(), common.UserCredentials{
		Login:    user.Login,
		Password: user.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "userID",
		Value: newUser.ID,
	}

	if err = auth.WriteEncryptedCookie(w, cookie); err != nil {
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
