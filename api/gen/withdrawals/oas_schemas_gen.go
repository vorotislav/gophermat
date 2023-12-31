// Code generated by ogen, DO NOT EDIT.

package api

import (
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

// GetWithdrawalsInternalServerError is response for GetWithdrawals operation.
type GetWithdrawalsInternalServerError struct{}

func (*GetWithdrawalsInternalServerError) getWithdrawalsRes() {}

// GetWithdrawalsNoContent is response for GetWithdrawals operation.
type GetWithdrawalsNoContent struct{}

func (*GetWithdrawalsNoContent) getWithdrawalsRes() {}

type GetWithdrawalsOKApplicationJSON []GetWithdrawalsOKItem

func (*GetWithdrawalsOKApplicationJSON) getWithdrawalsRes() {}

type GetWithdrawalsOKItem struct {
	Order       OptString   `json:"order"`
	Sum         OptFloat64  `json:"sum"`
	ProcessedAt OptDateTime `json:"processed_at"`
}

// GetOrder returns the value of Order.
func (s *GetWithdrawalsOKItem) GetOrder() OptString {
	return s.Order
}

// GetSum returns the value of Sum.
func (s *GetWithdrawalsOKItem) GetSum() OptFloat64 {
	return s.Sum
}

// GetProcessedAt returns the value of ProcessedAt.
func (s *GetWithdrawalsOKItem) GetProcessedAt() OptDateTime {
	return s.ProcessedAt
}

// SetOrder sets the value of Order.
func (s *GetWithdrawalsOKItem) SetOrder(val OptString) {
	s.Order = val
}

// SetSum sets the value of Sum.
func (s *GetWithdrawalsOKItem) SetSum(val OptFloat64) {
	s.Sum = val
}

// SetProcessedAt sets the value of ProcessedAt.
func (s *GetWithdrawalsOKItem) SetProcessedAt(val OptDateTime) {
	s.ProcessedAt = val
}

// GetWithdrawalsUnauthorized is response for GetWithdrawals operation.
type GetWithdrawalsUnauthorized struct{}

func (*GetWithdrawalsUnauthorized) getWithdrawalsRes() {}

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

// NewOptFloat64 returns new OptFloat64 with value set to v.
func NewOptFloat64(v float64) OptFloat64 {
	return OptFloat64{
		Value: v,
		Set:   true,
	}
}

// OptFloat64 is optional float64.
type OptFloat64 struct {
	Value float64
	Set   bool
}

// IsSet returns true if OptFloat64 was set.
func (o OptFloat64) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptFloat64) Reset() {
	var v float64
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptFloat64) SetTo(v float64) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptFloat64) Get() (v float64, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptFloat64) Or(d float64) float64 {
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
