package register

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-faster/errors"
	"gophermat/internal/models"
	"net/http"

	api "gophermat/api/gen/register"

	"go.uber.org/zap"
)

const (
	APIRegisterPath = "/register"
)

type gmart interface {
	RegisterUser(ctx context.Context, user models.User) (string, error)
}

type register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		h.log.Info(fmt.Sprintf("Failed to user register: unknown Content-Type: %s", contentType))

		http.Error(w, "unknown Content-Type", http.StatusBadRequest)

		return
	}

	reg := register{}
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		h.log.Info(fmt.Sprintf("Failed to user register: cannot decode register data: %s", err.Error()))

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := h.gmart.RegisterUser(r.Context(), models.User{
		Login:    reg.Login,
		Password: reg.Password,
	})

	if err != nil {
		h.log.Info(fmt.Sprintf("Failed to user register: %s", err.Error()))

		if errors.Is(err, models.ErrConflict) {
			http.Error(w, "a user with this login is already registered", http.StatusConflict)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RegisterUser(ctx context.Context, req api.OptRegisterUserReq) (api.RegisterUserRes, error) {
	token, err := h.gmart.RegisterUser(ctx, models.User{
		Login:    req.Value.Login,
		Password: req.Value.Password,
	})

	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return &api.RegisterUserConflict{}, nil
		}

		return &api.RegisterUserInternalServerError{}, nil
	}

	return &api.RegisterUserOK{
		Data: api.NewOptRegisterUserOKData(api.RegisterUserOKData{
			Token: api.NewOptString(token),
		})}, nil
}
