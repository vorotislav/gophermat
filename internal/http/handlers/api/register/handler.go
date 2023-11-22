package register

import (
	"context"
	"github.com/go-faster/errors"
	"gophermat/internal/models"

	api "gophermat/api/gen/register"

	"go.uber.org/zap"
)

const (
	APIRegisterPath = "/register"
)

type gmart interface {
	RegisterUser(ctx context.Context, user models.User) (string, error)
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
