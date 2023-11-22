// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// LoginUser implements loginUser operation.
//
// POST /api/user/login
func (UnimplementedHandler) LoginUser(ctx context.Context, req OptLoginUserReq) (r LoginUserRes, _ error) {
	return r, ht.ErrNotImplemented
}
