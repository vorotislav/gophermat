// Code generated by ogen, DO NOT EDIT.

package api

import (
	"io"
	"mime"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

func decodeDeductPointsResponse(resp *http.Response) (res DeductPointsRes, _ error) {
	switch resp.StatusCode {
	case 200:
		// Code 200.
		return &DeductPointsOK{}, nil
	case 401:
		// Code 401.
		return &DeductPointsUnauthorized{}, nil
	case 402:
		// Code 402.
		return &DeductPointsPaymentRequired{}, nil
	case 422:
		// Code 422.
		return &DeductPointsUnprocessableEntity{}, nil
	case 500:
		// Code 500.
		return &DeductPointsInternalServerError{}, nil
	}
	return res, validate.UnexpectedStatusCode(resp.StatusCode)
}

func decodeGetBalanceResponse(resp *http.Response) (res GetBalanceRes, _ error) {
	switch resp.StatusCode {
	case 200:
		// Code 200.
		ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			return res, errors.Wrap(err, "parse media type")
		}
		switch {
		case ct == "application/json":
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				return res, err
			}
			d := jx.DecodeBytes(buf)

			var response GetBalanceOK
			if err := func() error {
				if err := response.Decode(d); err != nil {
					return err
				}
				if err := d.Skip(); err != io.EOF {
					return errors.New("unexpected trailing data")
				}
				return nil
			}(); err != nil {
				err = &ogenerrors.DecodeBodyError{
					ContentType: ct,
					Body:        buf,
					Err:         err,
				}
				return res, err
			}
			// Validate response.
			if err := func() error {
				if err := response.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return res, errors.Wrap(err, "validate")
			}
			return &response, nil
		default:
			return res, validate.InvalidContentType(ct)
		}
	case 401:
		// Code 401.
		return &GetBalanceUnauthorized{}, nil
	case 500:
		// Code 500.
		return &GetBalanceInternalServerError{}, nil
	}
	return res, validate.UnexpectedStatusCode(resp.StatusCode)
}

func decodeGetWithdrawalsResponse(resp *http.Response) (res GetWithdrawalsRes, _ error) {
	switch resp.StatusCode {
	case 200:
		// Code 200.
		ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			return res, errors.Wrap(err, "parse media type")
		}
		switch {
		case ct == "application/json":
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				return res, err
			}
			d := jx.DecodeBytes(buf)

			var response GetWithdrawalsOKApplicationJSON
			if err := func() error {
				if err := response.Decode(d); err != nil {
					return err
				}
				if err := d.Skip(); err != io.EOF {
					return errors.New("unexpected trailing data")
				}
				return nil
			}(); err != nil {
				err = &ogenerrors.DecodeBodyError{
					ContentType: ct,
					Body:        buf,
					Err:         err,
				}
				return res, err
			}
			// Validate response.
			if err := func() error {
				if err := response.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return res, errors.Wrap(err, "validate")
			}
			return &response, nil
		default:
			return res, validate.InvalidContentType(ct)
		}
	case 204:
		// Code 204.
		return &GetWithdrawalsNoContent{}, nil
	case 401:
		// Code 401.
		return &GetWithdrawalsUnauthorized{}, nil
	case 500:
		// Code 500.
		return &GetWithdrawalsInternalServerError{}, nil
	}
	return res, validate.UnexpectedStatusCode(resp.StatusCode)
}