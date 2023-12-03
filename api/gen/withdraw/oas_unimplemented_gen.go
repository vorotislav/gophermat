// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// DeductPoints implements deductPoints operation.
//
// POST /api/user/balance/withdraw
func (UnimplementedHandler) DeductPoints(ctx context.Context, req OptDeductPointsReq) (r DeductPointsRes, _ error) {
	return r, ht.ErrNotImplemented
}
