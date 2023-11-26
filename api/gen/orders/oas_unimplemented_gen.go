// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// GetOrders implements getOrders operation.
//
// GET /api/user/orders
func (UnimplementedHandler) GetOrders(ctx context.Context) (r GetOrdersRes, _ error) {
	return r, ht.ErrNotImplemented
}

// LoadOrder implements loadOrder operation.
//
// POST /api/user/orders
func (UnimplementedHandler) LoadOrder(ctx context.Context, req LoadOrderReq) (r LoadOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}
