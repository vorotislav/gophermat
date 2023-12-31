package app

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gophermat/internal/crypt"
	"gophermat/internal/luhn"
	"strconv"
	"time"

	"gophermat/internal/models"
)

var (
	errTokenPayload = errors.New("cannot get token payload from context")
	errProcessing   = errors.New("processing order error")
)

const (
	tickerDuration = time.Second * 2
	maxWorkers     = 10
	maxCapacity    = 50
	newStatusOrder = "NEW"
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
	GetNotProcessOrders() ([]models.Order, error)
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
	doneCh  chan struct{}
	pool    *pond.WorkerPool
	eg      errgroup.Group
}

func NewGMart(log *zap.Logger, auth authorizer, storage storage, ac accrualClient) *GMart {
	gm := &GMart{
		log:     log,
		auth:    auth,
		storage: storage,
		client:  ac,
		doneCh:  make(chan struct{}),
		pool:    pond.New(maxWorkers, maxCapacity),
		eg:      errgroup.Group{},
	}

	gm.eg.Go(func() error {
		err := processingAccrualOrders(log, storage, ac, gm.doneCh, gm.pool)
		if err != nil {
			return fmt.Errorf("%w: %w", errProcessing, err)
		}

		return nil
	})

	return gm
}

func (gm *GMart) Stop() {
	gm.log.Info("GMart stop. close channel and pool")
	close(gm.doneCh)
	gm.pool.Stop()
	if err := gm.eg.Wait(); err != nil {
		gm.log.Error("cannot wait goroutine finish", zap.Error(err))
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
		gm.log.Error("order number is not correct on luhn", zap.String("order number", orderNumber))

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
		gm.log.Info("the order already exists", zap.String("number", orderNumber))

		if order.UserID == tokenPayload.UserID {
			return models.ErrOrderUploaded
		}

		return models.ErrOrderUploadedAnotherUser
	}

	o := models.Order{
		UserID:     tokenPayload.UserID,
		Number:     orderNumber,
		Status:     newStatusOrder,
		UploadedAt: time.Now(),
	}

	err = gm.storage.SaveOrder(ctx, o)
	if err != nil {
		gm.log.Error("cannot save order", zap.Error(err))

		return err
	}

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
		return models.TokenPayload{}, errTokenPayload
	}

	tokenPayload, ok := value.(models.TokenPayload)
	if !ok {
		return models.TokenPayload{}, errTokenPayload
	}

	return tokenPayload, nil
}

func processingAccrualOrders(
	log *zap.Logger,
	store storage,
	client accrualClient,
	doneCh chan struct{},
	pool *pond.WorkerPool) error {
	tick := time.NewTicker(tickerDuration)

	for {
		select {
		case <-doneCh:
			return nil
		case <-tick.C:
			{
				orders, err := store.GetNotProcessOrders()
				if err != nil {
					if errors.Is(err, models.ErrNotFound) {
						continue
					}
					log.Warn("cannot get not process orders", zap.Error(err))
					continue
				}

				for _, o := range orders {
					o := o
					pool.Submit(func() {
						processOrder(log, store, client, o)
					})
				}

			}

		}
	}
}

func processOrder(log *zap.Logger, store storage, client accrualClient, order models.Order) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second+5)
	defer cancel()

	accrual, err := client.GetOrderAccrual(ctx, order.Number)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			log.Debug("order not found in accrual", zap.String("order number", order.Number))

			return
		}

		log.Error("cannot get order accrual", zap.Error(err))

		return
	}

	err = store.UpdateOrder(ctx, order.Number, accrual.Status, int(accrual.Accrual*100))
	if err != nil {
		log.Error("cannot update order accrual", zap.Error(err))

		return
	}

	log.Info("order successful updated",
		zap.String("order number", order.Number),
		zap.String("status", accrual.Status),
		zap.Float32("accrual", accrual.Accrual))

	balance, err := store.GetBalance(ctx, order.UserID)
	if err != nil {
		log.Debug("cannot get balance", zap.Error(err))
	}

	err = store.UpdateBalance(ctx, models.Balance{
		Current:  balance.Current + int(accrual.Accrual*100),
		Withdraw: balance.Withdraw,
	}, order.UserID)
	if err != nil {
		log.Error("cannot update balance", zap.Error(err))

		return
	}
}
