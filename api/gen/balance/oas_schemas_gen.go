// Code generated by ogen, DO NOT EDIT.

package api

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

// DeductPointsInternalServerError is response for DeductPoints operation.
type DeductPointsInternalServerError struct{}

func (*DeductPointsInternalServerError) deductPointsRes() {}

// DeductPointsOK is response for DeductPoints operation.
type DeductPointsOK struct{}

func (*DeductPointsOK) deductPointsRes() {}

// DeductPointsPaymentRequired is response for DeductPoints operation.
type DeductPointsPaymentRequired struct{}

func (*DeductPointsPaymentRequired) deductPointsRes() {}

type DeductPointsReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// GetOrder returns the value of Order.
func (s *DeductPointsReq) GetOrder() string {
	return s.Order
}

// GetSum returns the value of Sum.
func (s *DeductPointsReq) GetSum() float64 {
	return s.Sum
}

// SetOrder sets the value of Order.
func (s *DeductPointsReq) SetOrder(val string) {
	s.Order = val
}

// SetSum sets the value of Sum.
func (s *DeductPointsReq) SetSum(val float64) {
	s.Sum = val
}

// DeductPointsUnauthorized is response for DeductPoints operation.
type DeductPointsUnauthorized struct{}

func (*DeductPointsUnauthorized) deductPointsRes() {}

// DeductPointsUnprocessableEntity is response for DeductPoints operation.
type DeductPointsUnprocessableEntity struct{}

func (*DeductPointsUnprocessableEntity) deductPointsRes() {}

// GetBalanceInternalServerError is response for GetBalance operation.
type GetBalanceInternalServerError struct{}

func (*GetBalanceInternalServerError) getBalanceRes() {}

type GetBalanceNoContent struct {
	Current   OptFloat64 `json:"current"`
	Withdrawn OptFloat64 `json:"withdrawn"`
}

// GetCurrent returns the value of Current.
func (s *GetBalanceNoContent) GetCurrent() OptFloat64 {
	return s.Current
}

// GetWithdrawn returns the value of Withdrawn.
func (s *GetBalanceNoContent) GetWithdrawn() OptFloat64 {
	return s.Withdrawn
}

// SetCurrent sets the value of Current.
func (s *GetBalanceNoContent) SetCurrent(val OptFloat64) {
	s.Current = val
}

// SetWithdrawn sets the value of Withdrawn.
func (s *GetBalanceNoContent) SetWithdrawn(val OptFloat64) {
	s.Withdrawn = val
}

func (*GetBalanceNoContent) getBalanceRes() {}

type GetBalanceOK struct {
	Current   OptFloat64 `json:"current"`
	Withdrawn OptFloat64 `json:"withdrawn"`
}

// GetCurrent returns the value of Current.
func (s *GetBalanceOK) GetCurrent() OptFloat64 {
	return s.Current
}

// GetWithdrawn returns the value of Withdrawn.
func (s *GetBalanceOK) GetWithdrawn() OptFloat64 {
	return s.Withdrawn
}

// SetCurrent sets the value of Current.
func (s *GetBalanceOK) SetCurrent(val OptFloat64) {
	s.Current = val
}

// SetWithdrawn sets the value of Withdrawn.
func (s *GetBalanceOK) SetWithdrawn(val OptFloat64) {
	s.Withdrawn = val
}

func (*GetBalanceOK) getBalanceRes() {}

// GetBalanceUnauthorized is response for GetBalance operation.
type GetBalanceUnauthorized struct{}

func (*GetBalanceUnauthorized) getBalanceRes() {}

// NewOptDeductPointsReq returns new OptDeductPointsReq with value set to v.
func NewOptDeductPointsReq(v DeductPointsReq) OptDeductPointsReq {
	return OptDeductPointsReq{
		Value: v,
		Set:   true,
	}
}

// OptDeductPointsReq is optional DeductPointsReq.
type OptDeductPointsReq struct {
	Value DeductPointsReq
	Set   bool
}

// IsSet returns true if OptDeductPointsReq was set.
func (o OptDeductPointsReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDeductPointsReq) Reset() {
	var v DeductPointsReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDeductPointsReq) SetTo(v DeductPointsReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDeductPointsReq) Get() (v DeductPointsReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDeductPointsReq) Or(d DeductPointsReq) DeductPointsReq {
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
