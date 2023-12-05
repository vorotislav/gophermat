// Code generated by ogen, DO NOT EDIT.

package api

import (
	"math/bits"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/validate"
)

// Encode implements json.Marshaler.
func (s *DeductPointsReq) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *DeductPointsReq) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("order")
		e.Str(s.Order)
	}
	{
		e.FieldStart("sum")
		e.Float64(s.Sum)
	}
}

var jsonFieldsNameOfDeductPointsReq = [2]string{
	0: "order",
	1: "sum",
}

// Decode decodes DeductPointsReq from json.
func (s *DeductPointsReq) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode DeductPointsReq to nil")
	}
	var requiredBitSet [1]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "order":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Str()
				s.Order = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"order\"")
			}
		case "sum":
			requiredBitSet[0] |= 1 << 1
			if err := func() error {
				v, err := d.Float64()
				s.Sum = float64(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"sum\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode DeductPointsReq")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [1]uint8{
		0b00000011,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfDeductPointsReq) {
					name = jsonFieldsNameOfDeductPointsReq[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *DeductPointsReq) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *DeductPointsReq) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *GetBalanceNoContent) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *GetBalanceNoContent) encodeFields(e *jx.Encoder) {
	{
		if s.Current.Set {
			e.FieldStart("current")
			s.Current.Encode(e)
		}
	}
	{
		if s.Withdrawn.Set {
			e.FieldStart("withdrawn")
			s.Withdrawn.Encode(e)
		}
	}
}

var jsonFieldsNameOfGetBalanceNoContent = [2]string{
	0: "current",
	1: "withdrawn",
}

// Decode decodes GetBalanceNoContent from json.
func (s *GetBalanceNoContent) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode GetBalanceNoContent to nil")
	}

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "current":
			if err := func() error {
				s.Current.Reset()
				if err := s.Current.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"current\"")
			}
		case "withdrawn":
			if err := func() error {
				s.Withdrawn.Reset()
				if err := s.Withdrawn.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"withdrawn\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode GetBalanceNoContent")
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *GetBalanceNoContent) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *GetBalanceNoContent) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *GetBalanceOK) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *GetBalanceOK) encodeFields(e *jx.Encoder) {
	{
		if s.Current.Set {
			e.FieldStart("current")
			s.Current.Encode(e)
		}
	}
	{
		if s.Withdrawn.Set {
			e.FieldStart("withdrawn")
			s.Withdrawn.Encode(e)
		}
	}
}

var jsonFieldsNameOfGetBalanceOK = [2]string{
	0: "current",
	1: "withdrawn",
}

// Decode decodes GetBalanceOK from json.
func (s *GetBalanceOK) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode GetBalanceOK to nil")
	}

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "current":
			if err := func() error {
				s.Current.Reset()
				if err := s.Current.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"current\"")
			}
		case "withdrawn":
			if err := func() error {
				s.Withdrawn.Reset()
				if err := s.Withdrawn.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"withdrawn\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode GetBalanceOK")
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *GetBalanceOK) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *GetBalanceOK) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode encodes DeductPointsReq as json.
func (o OptDeductPointsReq) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	o.Value.Encode(e)
}

// Decode decodes DeductPointsReq from json.
func (o *OptDeductPointsReq) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptDeductPointsReq to nil")
	}
	o.Set = true
	if err := o.Value.Decode(d); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptDeductPointsReq) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptDeductPointsReq) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode encodes float64 as json.
func (o OptFloat64) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	e.Float64(float64(o.Value))
}

// Decode decodes float64 from json.
func (o *OptFloat64) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptFloat64 to nil")
	}
	o.Set = true
	v, err := d.Float64()
	if err != nil {
		return err
	}
	o.Value = float64(v)
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptFloat64) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptFloat64) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}
