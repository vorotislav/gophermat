package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"go.uber.org/zap"
	"gophermat/internal/models"
	"io"
	"net/http"
	"time"
)

const (
	httpClientTimeout = time.Millisecond * 500
)

type Client struct {
	dc        *http.Client
	log       *zap.Logger
	serverURL string
}

func NewClient(log *zap.Logger, accrualAddress string) *Client {
	c := &Client{
		dc: &http.Client{
			Timeout: httpClientTimeout,
		},
		log:       log,
		serverURL: fmt.Sprintf("%s/api/orders/", accrualAddress),
	}

	log.Debug("Client for accrual server", zap.String("url", c.serverURL))

	return c
}

func (c *Client) GetOrderAccrual(ctx context.Context, orderNumber string) (models.OrderAccrual, error) {
	c.log.Debug("new request for order", zap.String("order number", orderNumber))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.serverURL+"/"+orderNumber, http.NoBody)
	if err != nil {
		c.log.Error("cannot prepare request", zap.Error(err))

		return models.OrderAccrual{}, fmt.Errorf("cannot prepare request: %w", err)
	}

	var body []byte

	err = retry.Do(
		func() error {
			resp, err := c.dc.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil || resp.StatusCode >= http.StatusInternalServerError {
				return err
			}

			return nil
		},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.Attempts(4),
		retry.Context(ctx))

	var accrual models.OrderAccrual

	if err := json.Unmarshal(body, &accrual); err != nil {
		c.log.Error("cannot decode response", zap.Error(err))

		return models.OrderAccrual{}, fmt.Errorf("cannot decode accrual: %w", err)
	}

	return accrual, nil
}
