package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gophermat/internal/models"
	"gophermat/internal/settings"
	"net/http"
	"time"
)

type Client struct {
	dc        *http.Client
	log       *zap.Logger
	set       *settings.Settings
	serverURL string
}

func NewClient(log *zap.Logger, set *settings.Settings) *Client {
	c := &Client{
		dc: &http.Client{
			Timeout: time.Millisecond * 500,
		},
		log:       log,
		set:       set,
		serverURL: fmt.Sprintf("http://%s/api/orders/", set.AccrualSystemAddress),
	}

	return c
}

func (c *Client) GetOrderAccrual(ctx context.Context, orderNumber string) (models.OrderAccrual, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.serverURL+"/"+orderNumber, http.NoBody)
	if err != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot prepare request: %w", err)
	}

	resp, err := c.dc.Do(req)
	if err != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot get accrual: %w", err)
	}

	defer resp.Body.Close()

	var accrual models.OrderAccrual

	if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot decode accrual: %w", err)
	}

	return accrual, nil
}
