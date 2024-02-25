package api_error

import (
	"net/http"
)

type ApiError struct {
	Code        int
	MessageCode string
	Err         error
}

var _ error = (*ApiError)(nil)

func (e ApiError) Error() string {
	return e.Err.Error()
}

func (e ApiError) Unwrap() error {
	return e.Err
}

func NewApiError(statuscode int, msgCode string, err error) ApiError {
	return ApiError{
		Code:        statuscode,
		Err:         err,
		MessageCode: msgCode,
	}
}

func NewBadRequestError(code string, err error) ApiError {
	if code == "" {
		code = "something_went_wrong"
	}
	return NewApiError(http.StatusBadRequest, code, err)
}

func NewUnauthorizedError(code string, err error) ApiError {
	if code == "" {
		code = "missing_authentication_data"
	}
	return NewApiError(http.StatusUnauthorized, code, err)
}

func NewForbiddenError(code string, err error) ApiError {
	if code == "" {
		code = "not_allowed_to_perform_this_request"
	}
	return NewApiError(http.StatusForbidden, code, err)
}

func NewNotFoundError(code string, err error) ApiError {
	if code == "" {
		code = "resource_not_found"
	}
	return NewApiError(http.StatusNotFound, code, err)
}
