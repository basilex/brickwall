package common

import (
	"errors"
	"net/http"
)

var (
	// common
	ErrNotImplemented = errors.New("not implemented")

	// Request layer errors
	ErrReqBindJson = errors.New("failed to bind json")

	// Database layer errors
	ErrDBConnPoolExhausted = errors.New("failed to acquire connection")

	ErrDBNotFound    = errors.New("resource not found")
	ErrDBTrxError    = errors.New("transaction error")
	ErrDBTooManyRows = errors.New("too many rows")

	ErrDBRecordCount  = errors.New("failed to count record(s)")
	ErrDBRecordInsert = errors.New("failed to insert record(s)")
	ErrDBRecordSelect = errors.New("failed to select record(s)")
	ErrDBRecordUpdate = errors.New("failed to update record(s)")
	ErrDBRecordDelete = errors.New("failed to delete record(s)")

	// Auth/jwt layer errors
	ErrJwtTokenInvalid     = errors.New("invalid token")
	ErrJwtTokenSigning     = errors.New("signing token")
	ErrJwtTokenClaims      = errors.New("invalid claims")
	ErrJwtTokenInvalidated = errors.New("invalidated token")

	// 2FA layer errors
	Err2FAKeyGeneration = errors.New("failed to generate key")

	// Business layer errors

	// Network layer errors

	// Other layer errors
)

type Exception struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
}

func NewException(code int, message string) *Exception {
	return &Exception{
		Code:      code,
		Message:   message,
		TimeStamp: TimestampISO3339NS(),
	}
}

func ErrMapper(err error) (int, *Exception) {
	switch {
	case errors.Is(err, ErrReqBindJson):
		return http.StatusBadRequest, NewException(http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrDBConnPoolExhausted):
		return http.StatusServiceUnavailable, NewException(http.StatusServiceUnavailable, err.Error())

	case errors.Is(err, ErrDBNotFound):
		return http.StatusNotFound, NewException(http.StatusNotFound, err.Error())
	case errors.Is(err, ErrDBTrxError):
		return http.StatusInternalServerError, NewException(http.StatusInternalServerError, err.Error())
	case errors.Is(err, ErrDBTooManyRows):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())

	case errors.Is(err, ErrDBRecordCount):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())
	case errors.Is(err, ErrDBRecordInsert):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())
	case errors.Is(err, ErrDBRecordSelect):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())
	case errors.Is(err, ErrDBRecordUpdate):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())
	case errors.Is(err, ErrDBRecordDelete):
		return http.StatusNotAcceptable, NewException(http.StatusNotAcceptable, err.Error())

	case errors.Is(err, ErrJwtTokenInvalid):
		return http.StatusUnauthorized, NewException(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrJwtTokenSigning):
		return http.StatusUnauthorized, NewException(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrJwtTokenInvalidated):
		return http.StatusUnauthorized, NewException(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrJwtTokenClaims):
		return http.StatusUnauthorized, NewException(http.StatusUnauthorized, err.Error())

	case errors.Is(err, Err2FAKeyGeneration):
		return http.StatusExpectationFailed, NewException(http.StatusExpectationFailed, err.Error())
	default:
		return http.StatusInternalServerError, NewException(http.StatusInternalServerError, err.Error())
	}
}
