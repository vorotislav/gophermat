// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// GetWithdrawals implements getWithdrawals operation.
//
// GET /api/user/withdrawals
func (UnimplementedHandler) GetWithdrawals(ctx context.Context) (r GetWithdrawalsRes, _ error) {
	return r, ht.ErrNotImplemented
}
