package login

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"gophermat/internal/models"
	"net/http"

	api "gophermat/api/gen/login"
)

const (
	APILoginPath = "/login"
)

type login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type gmart interface {
	LoginUser(ctx context.Context, user models.User) (string, error)
}

type Handler struct {
	log *zap.Logger

	gmart gmart
}

func NewHandler(log *zap.Logger, gmart gmart) *Handler {
	return &Handler{
		log:   log,
		gmart: gmart,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		h.log.Info(fmt.Sprintf("Failed to user login: unknown Content-Type: %s", contentType))

		http.Error(w, "unknown Content-Type", http.StatusBadRequest)

		return
	}

	l := login{}
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		h.log.Info(fmt.Sprintf("Failed to user login: cannot decode login data: %s", err.Error()))

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := h.gmart.LoginUser(r.Context(), models.User{
		Login:    l.Login,
		Password: l.Password,
	})

	if err != nil {
		h.log.Info(fmt.Sprintf("Failed to user login: %s", err.Error()))

		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "user not found", http.StatusUnauthorized)

			return
		}

		if errors.Is(err, models.ErrInvalidPassword) {
			http.Error(w, "invalid password", http.StatusUnauthorized)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(ctx context.Context, req api.OptLoginUserReq) (api.LoginUserRes, error) {
	token, err := h.gmart.LoginUser(ctx, models.User{
		Login:    req.Value.Login,
		Password: req.Value.Password,
	})

	if err != nil {
		if errors.Is(err, models.ErrNotFound) ||
			errors.Is(err, models.ErrInvalidPassword) {
			return &api.LoginUserUnauthorized{}, nil
		}

		return &api.LoginUserInternalServerError{}, nil
	}

	return &api.LoginUserOK{
		Data: api.NewOptLoginUserOKData(api.LoginUserOKData{
			Token: api.NewOptString(token),
		}),
	}, nil
}
