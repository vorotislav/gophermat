package app

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"gophermat/internal/crypt"
	"gophermat/internal/luhn"
	"strconv"
	"time"

	"gophermat/internal/models"
)

var (
	ErrTokenPayload = errors.New("cannot get token payload from context")
)

type storage interface {
	RegisterUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, user models.User) (models.User, error)
	GetOrder(ctx context.Context, orderNumber string) (models.Order, error)
	SaveOrder(ctx context.Context, order models.Order) error
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
	UpdateOrder(ctx context.Context, orderNumber, status string, accrual int) error
}

type authorizer interface {
	GenerateToken(payload models.TokenPayload) (string, error)
}

type accrualClient interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (models.OrderAccrual, error)
}

type GMart struct {
	log     *zap.Logger
	auth    authorizer
	storage storage
	client  accrualClient
}

func NewGMart(log *zap.Logger, auth authorizer, storage storage, ac accrualClient) *GMart {
	return &GMart{
		log:     log,
		auth:    auth,
		storage: storage,
		client:  ac,
	}
}

func (gm *GMart) RegisterUser(ctx context.Context, user models.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInvalidInput, err)
	}

	_, err := gm.storage.GetUser(ctx, user)
	if err == nil {
		return "", models.ErrConflict
	}

	passHash, err := crypt.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	user.Password = passHash

	u, err := gm.storage.RegisterUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	token, err := gm.auth.GenerateToken(models.TokenPayload{UserID: u.ID})
	if err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	return token, nil
}

func (gm *GMart) LoginUser(ctx context.Context, user models.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInvalidInput, err)
	}

	u, err := gm.storage.GetUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return "", err
		}

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	if err := crypt.CheckPassword(user.Password, u.Password); err != nil {
		return "", models.ErrInvalidPassword
	}

	token, err := gm.auth.GenerateToken(models.TokenPayload{UserID: u.ID})
	if err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	return token, nil
}

func (gm *GMart) LoadOrder(ctx context.Context, orderNumber string) error {
	// проверяем корректность номера заказа
	on, err := strconv.Atoi(orderNumber)
	if err != nil {
		return fmt.Errorf("%w: %w", models.ErrInvalidOrderNumber, err)
	}

	// проверяем номер заказа по алгоритму Луна
	if ok := luhn.Valid(on); !ok {
		return fmt.Errorf("%w: order number is not correct on luhn", models.ErrInvalidOrderNumber)
	}

	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		return err
	}

	// проверяем номер заказа в репозитории
	order, err := gm.storage.GetOrder(ctx, orderNumber)
	if err == nil {
		if order.UserID == tokenPayload.UserID {
			return models.ErrOrderUploaded
		}

		return models.ErrOrderUploadedAnotherUser
	}

	o := models.Order{
		UserID:     tokenPayload.UserID,
		Number:     orderNumber,
		UploadedAt: time.Now(),
	}

	err = gm.storage.SaveOrder(ctx, o)
	if err != nil {
		return err
	}

	go gm.getOrderAccrual(orderNumber)

	return nil
}

func (gm *GMart) GetOrders(ctx context.Context) ([]models.Order, error) {
	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	orders, err := gm.storage.GetOrders(ctx, tokenPayload.UserID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (gm *GMart) getOrderAccrual(orderNumber string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	accrual, err := gm.client.GetOrderAccrual(ctx, orderNumber)
	if err != nil {
		gm.log.Error("cannot get order accrual", zap.Error(err))

		return
	}

	err = gm.storage.UpdateOrder(ctx, orderNumber, accrual.Status, accrual.Accrual)
	if err != nil {
		gm.log.Error("cannot update order accrual", zap.Error(err))

		return
	}

	gm.log.Info("order successful updated", zap.String("order number", orderNumber))
}

func payloadFromContext(ctx context.Context) (models.TokenPayload, error) {
	value := ctx.Value(models.CtxTokenPayload{})
	if value == nil {
		return models.TokenPayload{}, ErrTokenPayload
	}

	tokenPayload, ok := value.(models.TokenPayload)
	if !ok {
		return models.TokenPayload{}, ErrTokenPayload
	}

	return tokenPayload, nil
}
