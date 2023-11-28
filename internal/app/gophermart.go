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
	GetBalance(ctx context.Context, userID int) (models.Balance, error)
	UpdateBalance(ctx context.Context, balance models.Balance, userID int) error
	AddBalanceHistory(ctx context.Context, orderNumber string, sum, userID int) error
	GetBalanceHistory(ctx context.Context, userID int) ([]models.BalanceWithdrawal, error)
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
		gm.log.Error("cannot input validate", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInvalidInput, err)
	}

	_, err := gm.storage.GetUser(ctx, user)
	if err == nil {
		gm.log.Info("this user is already registered")

		return "", models.ErrConflict
	} else {
		if !errors.Is(err, models.ErrNotFound) {
			gm.log.Error("cannot user register", zap.Error(err))

			return "", models.ErrInternal
		}
	}

	passHash, err := crypt.HashPassword(user.Password)
	if err != nil {
		gm.log.Error("cannot hash password", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	user.Password = passHash

	u, err := gm.storage.RegisterUser(ctx, user)
	if err != nil {
		gm.log.Error("cannot user register", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	token, err := gm.auth.GenerateToken(models.TokenPayload{UserID: u.ID})
	if err != nil {
		gm.log.Error("cannot generate token", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	return token, nil
}

func (gm *GMart) LoginUser(ctx context.Context, user models.User) (string, error) {
	if err := user.Validate(); err != nil {
		gm.log.Error("cannot input validate", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInvalidInput, err)
	}

	u, err := gm.storage.GetUser(ctx, user)
	if err != nil {
		gm.log.Error("cannot user login", zap.Error(err))

		if errors.Is(err, models.ErrNotFound) {
			return "", err
		}

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	if err := crypt.CheckPassword(user.Password, u.Password); err != nil {
		gm.log.Error("cannot user register", zap.Error(err))

		return "", models.ErrInvalidPassword
	}

	token, err := gm.auth.GenerateToken(models.TokenPayload{UserID: u.ID})
	if err != nil {
		gm.log.Error("cannot generate token", zap.Error(err))

		return "", fmt.Errorf("%w: %w", models.ErrInternal, err)
	}

	return token, nil
}

func (gm *GMart) LoadOrder(ctx context.Context, orderNumber string) error {
	// проверяем корректность номера заказа
	on, err := strconv.Atoi(orderNumber)
	if err != nil {
		gm.log.Error("cannot order number check", zap.Error(err))

		return fmt.Errorf("%w: %w", models.ErrInvalidOrderNumber, err)
	}

	// проверяем номер заказа по алгоритму Луна
	if ok := luhn.Valid(on); !ok {
		gm.log.Error("order number is not correct on luhn")

		return fmt.Errorf("%w: order number is not correct on luhn", models.ErrInvalidOrderNumber)
	}

	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		gm.log.Error("cannot get payload", zap.Error(err))

		return err
	}

	// проверяем номер заказа в репозитории
	order, err := gm.storage.GetOrder(ctx, orderNumber)
	if err == nil {
		gm.log.Error("cannot get order", zap.Error(err))

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
		gm.log.Error("cannot save order", zap.Error(err))

		return err
	}

	go gm.getOrderAccrual(orderNumber, tokenPayload.UserID)

	return nil
}

func (gm *GMart) GetOrders(ctx context.Context) ([]models.Order, error) {
	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		gm.log.Error("cannot get payload", zap.Error(err))

		return nil, err
	}

	orders, err := gm.storage.GetOrders(ctx, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot get orders", zap.Error(err))

		return nil, err
	}

	return orders, nil
}

func (gm *GMart) getOrderAccrual(orderNumber string, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second+60)
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

	if accrual.Accrual == 0 {
		return
	}

	balance, err := gm.storage.GetBalance(ctx, userID)
	if err != nil {
		gm.log.Error("cannot get balance", zap.Error(err))

		return
	}

	err = gm.storage.UpdateBalance(ctx, models.Balance{
		Current:  balance.Current + accrual.Accrual,
		Withdraw: balance.Withdraw,
	}, userID)
	if err != nil {
		gm.log.Error("cannot update balance", zap.Error(err))

		return
	}
}

func (gm *GMart) GetBalance(ctx context.Context) (models.Balance, error) {
	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		gm.log.Error("cannot get payload", zap.Error(err))

		return models.Balance{}, err
	}

	balance, err := gm.storage.GetBalance(ctx, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot get balance", zap.Error(err))

		return models.Balance{}, err
	}

	return balance, nil
}

func (gm *GMart) DeductPoints(ctx context.Context, withdraw models.BalanceWithdraw) error {
	// получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		gm.log.Error("cannot get payload", zap.Error(err))

		return err
	}

	// получаем баланс пользователя
	balance, err := gm.storage.GetBalance(ctx, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot get balance", zap.Error(err))

		return err
	}

	// проверяем, что текущий баланс позволяет списать запрошенную сумму
	if balance.Current < withdraw.Sum {
		gm.log.Error("the user's balance is less than the requested amount")

		return models.ErrInsufficientBalance
	}

	// меняем баланс
	err = gm.storage.UpdateBalance(ctx, models.Balance{
		Current:  balance.Current - withdraw.Sum,
		Withdraw: balance.Withdraw + withdraw.Sum,
	}, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot update balance", zap.Error(err))

		return err
	}

	// записываем изменение баланса в историю
	err = gm.storage.AddBalanceHistory(ctx, withdraw.Order, withdraw.Sum, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot add balance history", zap.Error(err))

		return err
	}

	return nil
}
func (gm *GMart) GetWithdrawals(ctx context.Context) ([]models.BalanceWithdrawal, error) { // получаем id пользователя
	tokenPayload, err := payloadFromContext(ctx)
	if err != nil {
		gm.log.Error("cannot get payload", zap.Error(err))

		return nil, err
	}

	history, err := gm.storage.GetBalanceHistory(ctx, tokenPayload.UserID)
	if err != nil {
		gm.log.Error("cannot get balance history", zap.Error(err))

		return nil, err
	}

	return history, nil
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
