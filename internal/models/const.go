package models

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrInvalidInput             = errors.New("invalid input")
	ErrConflict                 = errors.New("conflict")
	ErrInternal                 = errors.New("internal error")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrInvalidOrderNumber       = errors.New("invalid order number")
	ErrOrderUploaded            = errors.New("the order has already been uploaded")
	ErrOrderUploadedAnotherUser = errors.New("the order has already been uploaded another user")
	ErrInsufficientBalance      = errors.New("insufficient funds on the balance sheet")
)
