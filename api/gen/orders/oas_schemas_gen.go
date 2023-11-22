// Code generated by ogen, DO NOT EDIT.

package api

import (
	"io"
	"time"
)

type BearerAuth struct {
	Token string
}

// GetToken returns the value of Token.
func (s *BearerAuth) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *BearerAuth) SetToken(val string) {
	s.Token = val
}

type GetOrdersOKItem struct {
	Number     OptString   `json:"number"`
	Status     OptString   `json:"status"`
	Accrual    OptInt      `json:"accrual"`
	UploadedAt OptDateTime `json:"uploaded_at"`
}

// GetNumber returns the value of Number.
func (s *GetOrdersOKItem) GetNumber() OptString {
	return s.Number
}

// GetStatus returns the value of Status.
func (s *GetOrdersOKItem) GetStatus() OptString {
	return s.Status
}

// GetAccrual returns the value of Accrual.
func (s *GetOrdersOKItem) GetAccrual() OptInt {
	return s.Accrual
}

// GetUploadedAt returns the value of UploadedAt.
func (s *GetOrdersOKItem) GetUploadedAt() OptDateTime {
	return s.UploadedAt
}

// SetNumber sets the value of Number.
func (s *GetOrdersOKItem) SetNumber(val OptString) {
	s.Number = val
}

// SetStatus sets the value of Status.
func (s *GetOrdersOKItem) SetStatus(val OptString) {
	s.Status = val
}

// SetAccrual sets the value of Accrual.
func (s *GetOrdersOKItem) SetAccrual(val OptInt) {
	s.Accrual = val
}

// SetUploadedAt sets the value of UploadedAt.
func (s *GetOrdersOKItem) SetUploadedAt(val OptDateTime) {
	s.UploadedAt = val
}

// LoadOrderAccepted is response for LoadOrder operation.
type LoadOrderAccepted struct{}

func (*LoadOrderAccepted) loadOrderRes() {}

// LoadOrderBadRequest is response for LoadOrder operation.
type LoadOrderBadRequest struct{}

func (*LoadOrderBadRequest) loadOrderRes() {}

// LoadOrderConflict is response for LoadOrder operation.
type LoadOrderConflict struct{}

func (*LoadOrderConflict) loadOrderRes() {}

// LoadOrderInternalServerError is response for LoadOrder operation.
type LoadOrderInternalServerError struct{}

func (*LoadOrderInternalServerError) loadOrderRes() {}

// LoadOrderOK is response for LoadOrder operation.
type LoadOrderOK struct{}

func (*LoadOrderOK) loadOrderRes() {}

type LoadOrderReq struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s LoadOrderReq) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

// LoadOrderUnauthorized is response for LoadOrder operation.
type LoadOrderUnauthorized struct{}

func (*LoadOrderUnauthorized) loadOrderRes() {}

// LoadOrderUnprocessableEntity is response for LoadOrder operation.
type LoadOrderUnprocessableEntity struct{}

func (*LoadOrderUnprocessableEntity) loadOrderRes() {}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt returns new OptInt with value set to v.
func NewOptInt(v int) OptInt {
	return OptInt{
		Value: v,
		Set:   true,
	}
}

// OptInt is optional int.
type OptInt struct {
	Value int
	Set   bool
}

// IsSet returns true if OptInt was set.
func (o OptInt) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt) Reset() {
	var v int
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt) SetTo(v int) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt) Get() (v int, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}
