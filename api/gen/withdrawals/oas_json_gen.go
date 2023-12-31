// Code generated by ogen, DO NOT EDIT.

package api

import (
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/json"
)

// Encode encodes GetWithdrawalsOKApplicationJSON as json.
func (s GetWithdrawalsOKApplicationJSON) Encode(e *jx.Encoder) {
	unwrapped := []GetWithdrawalsOKItem(s)

	e.ArrStart()
	for _, elem := range unwrapped {
		elem.Encode(e)
	}
	e.ArrEnd()
}

// Decode decodes GetWithdrawalsOKApplicationJSON from json.
func (s *GetWithdrawalsOKApplicationJSON) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode GetWithdrawalsOKApplicationJSON to nil")
	}
	var unwrapped []GetWithdrawalsOKItem
	if err := func() error {
		unwrapped = make([]GetWithdrawalsOKItem, 0)
		if err := d.Arr(func(d *jx.Decoder) error {
			var elem GetWithdrawalsOKItem
			if err := elem.Decode(d); err != nil {
				return err
			}
			unwrapped = append(unwrapped, elem)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return errors.Wrap(err, "alias")
	}
	*s = GetWithdrawalsOKApplicationJSON(unwrapped)
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s GetWithdrawalsOKApplicationJSON) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *GetWithdrawalsOKApplicationJSON) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *GetWithdrawalsOKItem) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *GetWithdrawalsOKItem) encodeFields(e *jx.Encoder) {
	{
		if s.Order.Set {
			e.FieldStart("order")
			s.Order.Encode(e)
		}
	}
	{
		if s.Sum.Set {
			e.FieldStart("sum")
			s.Sum.Encode(e)
		}
	}
	{
		if s.ProcessedAt.Set {
			e.FieldStart("processed_at")
			s.ProcessedAt.Encode(e, json.EncodeDateTime)
		}
	}
}

var jsonFieldsNameOfGetWithdrawalsOKItem = [3]string{
	0: "order",
	1: "sum",
	2: "processed_at",
}

// Decode decodes GetWithdrawalsOKItem from json.
func (s *GetWithdrawalsOKItem) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode GetWithdrawalsOKItem to nil")
	}

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "order":
			if err := func() error {
				s.Order.Reset()
				if err := s.Order.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"order\"")
			}
		case "sum":
			if err := func() error {
				s.Sum.Reset()
				if err := s.Sum.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"sum\"")
			}
		case "processed_at":
			if err := func() error {
				s.ProcessedAt.Reset()
				if err := s.ProcessedAt.Decode(d, json.DecodeDateTime); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"processed_at\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode GetWithdrawalsOKItem")
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *GetWithdrawalsOKItem) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *GetWithdrawalsOKItem) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode encodes time.Time as json.
func (o OptDateTime) Encode(e *jx.Encoder, format func(*jx.Encoder, time.Time)) {
	if !o.Set {
		return
	}
	format(e, o.Value)
}

// Decode decodes time.Time from json.
func (o *OptDateTime) Decode(d *jx.Decoder, format func(*jx.Decoder) (time.Time, error)) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptDateTime to nil")
	}
	o.Set = true
	v, err := format(d)
	if err != nil {
		return err
	}
	o.Value = v
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptDateTime) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e, json.EncodeDateTime)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptDateTime) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d, json.DecodeDateTime)
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

// Encode encodes string as json.
func (o OptString) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	e.Str(string(o.Value))
}

// Decode decodes string from json.
func (o *OptString) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptString to nil")
	}
	o.Set = true
	v, err := d.Str()
	if err != nil {
		return err
	}
	o.Value = string(v)
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptString) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptString) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}
