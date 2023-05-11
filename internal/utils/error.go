package utils

import "net/http"

var (
	ErrInvalidPwd   = &APIError{Status: http.StatusBadRequest, msg: "password must contain 8 letters, 1 number, 1 upper case letter, 1 special character"}
	ErrInvalidEmail = &APIError{Status: http.StatusBadRequest, msg: "email must be valid"}
	ErrAppLayer     = &APIError{Status: http.StatusInternalServerError, msg: "unknown app error"}
	ErrDBLayer      = &APIError{Status: http.StatusInternalServerError, msg: "unknown storage error"}
	ErrThirdParty   = &APIError{Status: http.StatusInternalServerError, msg: "unknown third party service error"}
	ErrBadRequest   = &APIError{Status: http.StatusBadRequest, msg: "bad request"}
	ErrConflict     = &APIError{Status: http.StatusConflict, msg: "conflict"}

	ErrAuth           = &APIError{Status: http.StatusUnauthorized, msg: "invalid auth token"}
	ErrNotFound       = &APIError{Status: http.StatusNotFound, msg: "not found"}
	ErrDuplicate      = &APIError{Status: http.StatusConflict, msg: "duplicate"}
	ErrNotAuthorized  = &APIError{Status: http.StatusUnauthorized, msg: "not authorized"}
	ErrAlreadyCreated = &APIError{Status: http.StatusOK, msg: "entity already created"}
	ErrWrongFormat    = &APIError{Status: http.StatusUnprocessableEntity, msg: "entity provided has unproccessable format"}
	ErrNoData         = &APIError{Status: http.StatusNoContent, msg: "no data"}
	ErrPaymentError   = &APIError{Status: http.StatusPaymentRequired, msg: "not enough money to spend"}
)

type APIError struct {
	Status int
	msg    string
}

func (e APIError) Error() string {
	return e.msg
}

func (e APIError) APIError() (int, string) {
	return e.Status, e.msg
}

type WrappedAPIError struct {
	error
	APIError *APIError
}

func (we WrappedAPIError) Is(err error) bool {
	return we.APIError == err
}

func WrapError(err error, APIError *APIError) error {
	return WrappedAPIError{error: err, APIError: APIError}
}
