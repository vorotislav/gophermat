package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"gophermat/internal/http/handlers/api/orders"
	"gophermat/internal/models"
	"net/http"
	"time"

	apiLogin "gophermat/api/gen/login"
	apiOrders "gophermat/api/gen/orders"
	apiRegister "gophermat/api/gen/register"
	"gophermat/internal/http/handlers/api/login"
	"gophermat/internal/http/handlers/api/register"
	"gophermat/internal/settings"

	"github.com/go-chi/chi/v5"
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
	lr, err := apiLogin.NewServer(lh)
	if err != nil {
		return nil, err
	}

	routes = append(routes, Route{
		Pattern: APIPathPrefix + login.APILoginPath,
		Handler: lr,
	})

	rh := register.NewHandler(log, gmart)
	rr, err := apiRegister.NewServer(rh)
	if err != nil {
		return nil, err
	}

	routes = append(routes, Route{
		Pattern: APIPathPrefix + register.APIRegisterPath,
		Handler: rr,
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

	return routes, nil
}
