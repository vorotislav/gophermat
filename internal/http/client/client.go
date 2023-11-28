package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gophermat/internal/models"
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

	resp, err := c.dc.Do(req)
	if err != nil {
		c.log.Error("cannot do request", zap.Error(err))

		return models.OrderAccrual{}, fmt.Errorf("cannot get accrual: %w", err)
	}

	defer resp.Body.Close()

	var accrual models.OrderAccrual

	if resp.StatusCode != http.StatusOK {
		c.log.Info("response status not ok", zap.String("status", resp.Status))

		return models.OrderAccrual{}, fmt.Errorf("accrual request status %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
		c.log.Error("cannot decode response", zap.Error(err))

		return models.OrderAccrual{}, fmt.Errorf("cannot decode accrual: %w", err)
	}

	return accrual, nil
}
