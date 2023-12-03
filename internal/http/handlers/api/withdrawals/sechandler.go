package withdrawals

import (
	"context"
	"fmt"

	api "gophermat/api/gen/withdrawals"
	"gophermat/internal/models"
)

type authorizer interface {
	ParseToken(string) (models.TokenPayload, error)
}

type SecHandler struct {
	auth authorizer
}

func NewSecHandler(auth authorizer) *SecHandler {
	return &SecHandler{auth: auth}
}

func (s SecHandler) HandleBearerAuth(
	ctx context.Context,
	_ string,
	t api.BearerAuth,
) (context.Context, error) {
	tokenPayload, err := s.auth.ParseToken(t.Token)
	if err != nil {
		return ctx, fmt.Errorf("handled authorization: %w", err)
	}

	return context.WithValue(ctx, models.CtxTokenPayload{}, tokenPayload), nil
}
