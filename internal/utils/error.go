package utils

import "net/http"

var (
	ErrInvalidPwd   = &APIError{Status: http.StatusBadRequest, Msg: "password must contain 8 letters, 1 number, 1 upper case letter, 1 special character"}
	ErrInvalidEmail = &APIError{Status: http.StatusBadRequest, Msg: "email must be valid"}

	ErrAppLayer   = &APIError{Status: http.StatusInternalServerError, Msg: "unknown app error"}
	ErrDBLayer    = &APIError{Status: http.StatusInternalServerError, Msg: "unknown storage error"}
	ErrThirdParty = &APIError{Status: http.StatusInternalServerError, Msg: "unknown third party service error"}
	ErrBadRequest = &APIError{Status: http.StatusBadRequest, Msg: "bad request"}
	ErrConflict   = &APIError{Status: http.StatusConflict, Msg: "conflict"}

	ErrAuth           = &APIError{Status: http.StatusUnauthorized, Msg: "invalid auth token"}
	ErrNotFound       = &APIError{Status: http.StatusNotFound, Msg: "not found"}
	ErrDuplicate      = &APIError{Status: http.StatusConflict, Msg: "duplicate"}
	ErrNotAuthorized  = &APIError{Status: http.StatusUnauthorized, Msg: "not authorized"}
	ErrAlreadyCreated = &APIError{Status: http.StatusOK, Msg: "entity already created"}
	ErrWrongFormat    = &APIError{Status: http.StatusUnprocessableEntity, Msg: "entity provided has unproccessable format"}
	ErrNoData         = &APIError{Status: http.StatusNoContent, Msg: "no data"}
	ErrPaymentError   = &APIError{Status: http.StatusPaymentRequired, Msg: "not enough money to spend"}
)

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	return e.Msg
}

func (e APIError) APIError() (int, string) {
	return e.Status, e.Msg
}

type WrappedAPIError struct {
	error
	APIError *APIError
}

func (we WrappedAPIError) Is(err error) bool {
	return we.APIError == err
}

func (we WrappedAPIError) Message() string {
	return we.APIError.Msg
}

func (we WrappedAPIError) Unwrap() error {
	return we.error
}

func WrapError(err error, APIError *APIError) error {
	return WrappedAPIError{error: err, APIError: APIError}
}
