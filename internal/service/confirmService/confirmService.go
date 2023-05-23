package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type ConfirmationService struct {
	app *app.UserApp
}

func InitConfirmationService(a *app.UserApp) *ConfirmationService {
	return &ConfirmationService{a}
}

func (c *ConfirmationService) HandleConfirmUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	email := chi.URLParam(r, "email")
	confirmationCode := chi.URLParam(r, "confirmationCode")

	if email == "" || confirmationCode == "" {
		http.Error(w, "no email or code provided", http.StatusBadRequest)
		return
	}

	err := c.app.ConfirmUser(ctx, email, confirmationCode)

	if err != nil {
		switch err {
		case utils.ErrBadRequest:
			http.Error(w, "wrong code", http.StatusBadRequest)
			return
		case utils.ErrNotFound:
			http.Error(w, "no such user", http.StatusNoContent)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Account confirmed"))

	if err != nil {
		fmt.Print("error while confirming email: " + email)
	}
}
