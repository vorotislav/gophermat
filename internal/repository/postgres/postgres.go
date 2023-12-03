package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gophermat/internal/models"
	"time"
)

//go:embed migrations/*
var migrations embed.FS

const (
	storeDuration = time.Millisecond * 500
)

var (
	ErrSourceDriver   = errors.New("cannot create source driver")
	ErrSourceInstance = errors.New("cannot create migrate")
	ErrMigrateUp      = errors.New("cannot migrate up")
	ErrCreateStorage  = errors.New("cannot create storage")
)

type Storage struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, log *zap.Logger, databaseURI string) (*Storage, error) {
	log.Debug(fmt.Sprintf("Storage: database uri: %s", databaseURI))
	pool, err := pgxpool.New(ctx, databaseURI)

	if err != nil {
		log.Error("create pool", zap.Error(err))

		return nil, fmt.Errorf("%w: %w", ErrCreateStorage, err)
	}

	s := &Storage{
		log:  log,
		pool: pool,
	}

	if err = s.migrate(); err != nil {
		log.Error("migrations", zap.Error(err))

		return nil, fmt.Errorf("%w: %w", ErrCreateStorage, err)
	}

	return s, nil
}

func (s *Storage) migrate() error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("%w:%w", ErrSourceDriver, err)
	}

	connCfg := s.pool.Config().ConnConfig
	m, err := migrate.NewWithSourceInstance("iofs", d,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			connCfg.User, connCfg.Password, connCfg.Host, connCfg.Port, connCfg.Database))
	if err != nil {
		return fmt.Errorf("%w:%w", ErrSourceInstance, err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("%w:%w", ErrMigrateUp, err)
	}

	return nil
}

func (s *Storage) Stop() {
	s.pool.Close()
}

func (s *Storage) RegisterUser(ctx context.Context, user models.User) (models.User, error) {
	q := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"

	var id int

	err := s.pool.QueryRow(ctx, q, user.Login, user.Password).Scan(&id)
	if err != nil {
		return models.User{}, fmt.Errorf("cannot user register: %w", err)
	}

	user.ID = id

	return user, nil
}

func (s *Storage) GetUser(ctx context.Context, user models.User) (models.User, error) {
	q := "SELECT id, password FROM users WHERE login = $1"

	var (
		id       int
		password string
	)

	err := s.pool.QueryRow(ctx, q, user.Login).Scan(&id, &password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("cannot get user: %w", err)
	}

	user.ID = id
	user.Password = password

	return user, nil
}

func (s *Storage) GetOrder(ctx context.Context, orderNumber string) (models.Order, error) {
	q := "SELECT id, user_id, order_number, status, accrual, uploaded_at FROM orders WHERE order_number=$1"

	var o models.Order

	err := s.pool.QueryRow(ctx, q, orderNumber).Scan(
		&o.ID,
		&o.UserID,
		&o.Number,
		&o.Status,
		&o.Accrual,
		&o.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, models.ErrNotFound
		}

		return models.Order{}, fmt.Errorf("cannot get order by order number: %w", err)
	}

	return o, nil
}

func (s *Storage) SaveOrder(ctx context.Context, order models.Order) error {
	q := "INSERT INTO orders (user_id, order_number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := s.pool.Exec(ctx, q, order.UserID, order.Number, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		return fmt.Errorf("cannot save order: %w", err)
	}

	return nil
}

func (s *Storage) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	q := "SELECT id, user_id, order_number, status, accrual, uploaded_at FROM orders WHERE user_id=$1 ORDER BY uploaded_at"

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get orders: %w", err)
	}

	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		order := models.Order{}

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan orders: %w", err)
		}

		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	return orders, rows.Err()
}

func (s *Storage) GetNotProcessOrders() ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), storeDuration)
	defer cancel()

	q := "SELECT id, user_id, order_number, status FROM orders WHERE status not in ('INVALID', 'PROCESSED')"

	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("cannot get not processed orders")
	}

	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		order := models.Order{}

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status)
		if err != nil {
			return nil, fmt.Errorf("cannot scan orders: %w", err)
		}

		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("query err: %w", err)
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, orderNumber, status string, accrual int) error {
	q := "UPDATE orders SET(status, accrual) = ($1, $2) WHERE order_number = $3"

	_, err := s.pool.Exec(ctx, q, status, accrual, orderNumber)
	if err != nil {
		return fmt.Errorf("cannot update order: %w", err)
	}

	return nil
}

func (s *Storage) GetBalance(ctx context.Context, userID int) (models.Balance, error) {
	q := "SELECT current, withdraw FROM balance WHERE user_id = $1"

	b := models.Balance{}

	err := s.pool.QueryRow(ctx, q, userID).Scan(&b.Current, &b.Withdraw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Balance{}, models.ErrNotFound
		}

		return models.Balance{}, fmt.Errorf("cannot get balance: %w", err)
	}

	return b, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, balance models.Balance, userID int) error {
	q := "UPDATE balance SET(current, withdraw) = ($1, $2) WHERE user_id = $3"

	_, err := s.pool.Exec(ctx, q, balance.Current, balance.Withdraw, userID)
	if err != nil {
		return fmt.Errorf("cannot update balance: %w", err)
	}

	return nil
}

func (s *Storage) AddBalanceHistory(ctx context.Context, orderNumber string, sum, userID int) error {
	q := "INSERT INTO history (user_id, order_number, sum, processed_at) VALUES ($1, $2, $3, now())"

	_, err := s.pool.Exec(ctx, q, userID, orderNumber, sum)
	if err != nil {
		return fmt.Errorf("cannot insert balance history: %w", err)
	}

	return nil
}

func (s *Storage) GetBalanceHistory(ctx context.Context, userID int) ([]models.BalanceWithdrawal, error) {
	q := "SELECT order_number, sum, processed_at FROM history WHERE user_id = $1 ORDER BY processed_at"

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get balance history: %w", err)
	}

	defer rows.Close()

	history := make([]models.BalanceWithdrawal, 0)

	for rows.Next() {
		h := models.BalanceWithdrawal{}

		err = rows.Scan(&h.Order, &h.Sum, &h.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan balance history: %w", err)
		}

		history = append(history, h)
	}

	if len(history) == 0 {
		return nil, models.ErrNotFound
	}

	return history, rows.Err()
}
