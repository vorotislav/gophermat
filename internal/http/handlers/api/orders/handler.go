package orders

import (
	"context"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	api "gophermat/api/gen/orders"
	"gophermat/internal/models"
	"io"
)

const (
	APIOrdersPath = "/orders"
)

type gmart interface {
	LoadOrder(ctx context.Context, orderNumber string) error
	GetOrders(ctx context.Context) ([]models.Order, error)
}

type Handler struct {
	log   *zap.Logger
	gmart gmart
}

func NewHandler(log *zap.Logger, gmart gmart) *Handler {
	return &Handler{
		log:   log,
		gmart: gmart,
	}
}

func (h *Handler) GetOrders(ctx context.Context) (api.GetOrdersRes, error) {
	orders, err := h.gmart.GetOrders(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return &api.GetOrdersNoContent{}, nil
		}

		return &api.GetOrdersInternalServerError{}, err
	}

	result := make(api.GetOrdersOKApplicationJSON, 0, len(orders))
	for _, o := range orders {
		ro := api.GetOrdersOKItem{
			Number:     api.NewOptString(o.Number),
			Status:     api.NewOptString(o.Status),
			Accrual:    api.NewOptFloat64(float64(o.Accrual) / 100),
			UploadedAt: api.NewOptDateTime(o.UploadedAt),
		}

		result = append(result, ro)
	}

	return &result, nil
}

func (h *Handler) LoadOrder(ctx context.Context, req api.LoadOrderReq) (api.LoadOrderRes, error) {
	order, err := io.ReadAll(req.Data)
	if err != nil {
		return &api.LoadOrderInternalServerError{}, err
	}

	err = h.gmart.LoadOrder(ctx, string(order))
	if err != nil {
		if errors.Is(err, models.ErrInvalidOrderNumber) {
			return &api.LoadOrderUnprocessableEntity{}, nil
		}

		if errors.Is(err, models.ErrOrderUploaded) {
			return &api.LoadOrderOK{}, nil
		}

		if errors.Is(err, models.ErrOrderUploadedAnotherUser) {
			return &api.LoadOrderConflict{}, nil
		}

		return &api.LoadOrderInternalServerError{}, err
	}
	return &api.LoadOrderAccepted{}, nil
}
