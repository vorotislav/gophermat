package login

import (
	"context"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"gophermat/internal/models"

	api "gophermat/api/gen/login"
)

const (
	APILoginPath = "/login"
)

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
