package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidPwd   = &APIError{Code: codes.InvalidArgument, Status: http.StatusBadRequest, Msg: "password must contain 8 letters, 1 number, 1 upper case letter, 1 special character"}
	ErrInvalidEmail = &APIError{Code: codes.InvalidArgument, Status: http.StatusBadRequest, Msg: "email must be valid"}
	ErrInvalidTOTP  = &APIError{Code: codes.InvalidArgument, Status: http.StatusBadRequest, Msg: "totp must be valid"}

	ErrAppLayer   = &APIError{Code: codes.Internal, Status: http.StatusInternalServerError, Msg: "unknown app error"}
	ErrDBLayer    = &APIError{Code: codes.Internal, Status: http.StatusInternalServerError, Msg: "unknown storage error"}
	ErrThirdParty = &APIError{Code: codes.Internal, Status: http.StatusInternalServerError, Msg: "unknown third party service error"}
	ErrBadRequest = &APIError{Code: codes.InvalidArgument, Status: http.StatusBadRequest, Msg: "bad request"}
	ErrConflict   = &APIError{Code: codes.Unavailable, Status: http.StatusConflict, Msg: "conflict"}

	ErrNotFound      = &APIError{Code: codes.NotFound, Status: http.StatusNotFound, Msg: "not found"}
	ErrDuplicate     = &APIError{Code: codes.AlreadyExists, Status: http.StatusConflict, Msg: "duplicate"}
	ErrNotAuthorized = &APIError{Code: codes.InvalidArgument, Status: http.StatusUnauthorized, Msg: "not authorized"}
	ErrWrongFormat   = &APIError{Code: codes.InvalidArgument, Status: http.StatusUnprocessableEntity, Msg: "entity provided has unproccessable format"}
)

type APIError struct {
	Code   codes.Code
	Status int
	Msg    string
}

func (e APIError) Error() string {
	return e.Msg
}

func (e APIError) APIError() (code codes.Code, status int, msg string) {
	return e.Code, e.Status, e.Msg
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

func WrapError(err error, APIError *APIError) error {
	return WrappedAPIError{error: err, APIError: APIError}
}
