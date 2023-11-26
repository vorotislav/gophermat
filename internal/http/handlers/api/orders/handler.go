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

func (h *Handler) GetOrders(_ context.Context) (api.GetOrdersRes, error) {
	return nil, nil
}

func (h *Handler) LoadOrder(ctx context.Context, req api.LoadOrderReq) (api.LoadOrderRes, error) {
	order, err := io.ReadAll(req.Data)
	if err != nil {
		return &api.LoadOrderInternalServerError{}, err
	}

	err = h.gmart.LoadOrder(ctx, string(order))
	if err != nil {
		if errors.Is(err, models.ErrInvalidOrderNumber) {
			return &api.LoadOrderUnprocessableEntity{}, err
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
