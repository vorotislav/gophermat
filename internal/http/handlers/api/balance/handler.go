package balance

import (
	"context"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	api "gophermat/api/gen/balance"
	"gophermat/internal/models"
)

const (
	APIBalancePath = "/balance"
)

type gmart interface {
	GetBalance(ctx context.Context) (models.Balance, error)
	DeductPoints(ctx context.Context, withdraw models.BalanceWithdraw) error
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

func (h *Handler) DeductPoints(ctx context.Context, req api.OptDeductPointsReq) (api.DeductPointsRes, error) {
	err := h.gmart.DeductPoints(ctx, models.BalanceWithdraw{
		Order: req.Value.GetOrder(),
		Sum:   int(req.Value.GetSum() * 100),
	})

	if err != nil {
		if errors.Is(err, models.ErrInsufficientBalance) {
			return &api.DeductPointsPaymentRequired{}, nil
		}

		if errors.Is(err, models.ErrInvalidOrderNumber) {
			return &api.DeductPointsUnprocessableEntity{}, nil
		}

		return &api.DeductPointsInternalServerError{}, err
	}

	return &api.DeductPointsOK{}, nil
}

func (h *Handler) GetBalance(ctx context.Context) (api.GetBalanceRes, error) {
	balance, err := h.gmart.GetBalance(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &api.GetBalanceNoContent{}, nil
		}

		return &api.GetBalanceInternalServerError{}, err
	}

	return &api.GetBalanceOK{
		Current:   api.NewOptFloat64(float64(balance.Current) / 100),
		Withdrawn: api.NewOptInt(balance.Withdraw),
	}, nil
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
