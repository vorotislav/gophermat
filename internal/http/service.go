package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	apiBalance "gophermat/api/gen/balance"
	apiOrders "gophermat/api/gen/orders"
	apiWithdrawal "gophermat/api/gen/withdrawals"
	"gophermat/internal/http/handlers/api/balance"
	"gophermat/internal/http/handlers/api/login"
	"gophermat/internal/http/handlers/api/orders"
	"gophermat/internal/http/handlers/api/register"
	"gophermat/internal/http/handlers/api/withdrawals"
	"gophermat/internal/models"
	"gophermat/internal/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var (
	ErrCreateService = errors.New("create service")
)

const (
	APIPathPrefix = "/api/user"
)

type gmart interface {
	LoginUser(ctx context.Context, user models.User) (string, error)
	RegisterUser(ctx context.Context, user models.User) (string, error)
	LoadOrder(ctx context.Context, orderNumber string) error
	GetOrders(ctx context.Context) ([]models.Order, error)
	GetBalance(ctx context.Context) (models.Balance, error)
	DeductPoints(ctx context.Context, withdraw models.BalanceWithdraw) error
	GetWithdrawals(ctx context.Context) ([]models.BalanceWithdrawal, error)
}

type authorizer interface {
	ParseToken(string) (models.TokenPayload, error)
}

type Service struct {
	logger *zap.Logger
	server *http.Server

	gmart gmart
}

type Route struct {
	Pattern string
	Handler http.Handler
}

func NewService(log *zap.Logger, set *settings.Settings, gmart gmart, auth authorizer) (*Service, error) {
	mux := chi.NewRouter()

	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	rs, err := createRoutes(log, gmart, auth)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateService, err)
	}

	for _, route := range rs {
		mux.Mount(route.Pattern, route.Handler)
		log.Debug(fmt.Sprintf("added handler for %s", route.Pattern))
	}

	s := &http.Server{
		Addr:    set.Address,
		Handler: mux,
	}

	return &Service{
		logger: log.With(zap.String("package", "http service")),
		server: s,
		gmart:  gmart,
	}, nil
}

func (s *Service) Run() error {
	s.logger.Debug("Running server on", zap.String("address", s.server.Addr))

	return s.server.ListenAndServe()
}

func (s *Service) Stop(ctx context.Context) error {
	s.logger.Debug("stopping http service")

	return s.server.Shutdown(ctx)
}

func createRoutes(log *zap.Logger, gmart gmart, auth authorizer) ([]Route, error) {
	routes := make([]Route, 0)

	lh := login.NewHandler(log, gmart)

	routes = append(routes, Route{
		Pattern: APIPathPrefix + login.APILoginPath,
		Handler: http.HandlerFunc(lh.Login),
	})

	rh := register.NewHandler(log, gmart)

	routes = append(routes, Route{
		Pattern: APIPathPrefix + register.APIRegisterPath,
		Handler: http.HandlerFunc(rh.Register),
	})

	oh := orders.NewHandler(log, gmart)
	soh := orders.NewSecHandler(auth)
	or, err := apiOrders.NewServer(oh, soh)
	if err != nil {
		return nil, err
	}

	routes = append(routes, Route{
		Pattern: APIPathPrefix + orders.APIOrdersPath,
		Handler: or,
	})

	bh := balance.NewHandler(log, gmart)
	sbh := balance.NewSecHandler(auth)
	br, err := apiBalance.NewServer(bh, sbh)
	if err != nil {
		return nil, err
	}

	routes = append(routes, Route{
		Pattern: APIPathPrefix + balance.APIBalancePath,
		Handler: br,
	})

	wh := withdrawals.NewHandler(log, gmart)
	swh := withdrawals.NewSecHandler(auth)
	wr, err := apiWithdrawal.NewServer(wh, swh)
	if err != nil {
		return nil, err
	}

	routes = append(routes, Route{
		Pattern: APIPathPrefix + withdrawals.APIWithdrawalsPath,
		Handler: wr,
	})

	return routes, nil
}
