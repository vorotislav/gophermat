package withdrawals

import (
	"context"
	"errors"
	"go.uber.org/zap"

	api "gophermat/api/gen/withdrawals"
	"gophermat/internal/models"
)

const (
	APIWithdrawalsPath = "/withdrawals"
)

type gmart interface {
	GetWithdrawals(ctx context.Context) ([]models.BalanceWithdrawal, error)
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

func (h *Handler) GetWithdrawals(ctx context.Context) (api.GetWithdrawalsRes, error) {
	drawals, err := h.gmart.GetWithdrawals(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &api.GetWithdrawalsNoContent{}, nil
		}

		return &api.GetWithdrawalsInternalServerError{}, err
	}

	result := make(api.GetWithdrawalsOKApplicationJSON, 0, len(drawals))
	for _, d := range drawals {
		r := api.GetWithdrawalsOKItem{
			Order:       api.NewOptString(d.Order),
			Sum:         api.NewOptFloat64(float64(d.Sum) / 100),
			ProcessedAt: api.NewOptDateTime(d.ProcessedAt),
		}

		result = append(result, r)
	}

	return &result, nil
}
