package api

import (
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
)

var (
	// InvalidTokenErr error for an invalid Authorization header value.
	InvalidTokenErr = NewError(http.StatusUnauthorized, "Invalid token")
	// MissingTokenErr error for missing Authorization bearer token.
	MissingTokenErr = NewError(http.StatusBadRequest, "Missing bearer token")
	// InternalServerErr error for any internal server errors.
	InternalServerErr = NewError(http.StatusInternalServerError, "Internal server error")
)

// NewError constructor function for Error.
func NewError(code int, msg string) *Error {
	return &Error{code, msg}
}

// Error custom error type.
type Error struct {
	code    int
	Message string
}

// Error method to satisfy error interface.
func (e *Error) Error() string {
	return e.Message
}

// WriteError convenience function for writing error responses.
func WriteError(r *restful.Response, err interface{}) {
	var e *Error
	switch err.(type) {
	case Error, *Error:
		e = err.(*Error)
		break
	default:
		e = InternalServerErr
	}

	r.WriteHeaderAndEntity(e.code, e)
}

// BearerToken function returns the bearer token in an Authorization HTTP header.
func BearerToken(req *restful.Request) (string, error) {
	bearerPrefix := "Bearer "
	authField := req.HeaderParameter("Authorization")
	accessToken := strings.TrimPrefix(authField, bearerPrefix)

	if accessToken == "" || !strings.HasPrefix(authField, bearerPrefix) {
		return "", MissingTokenErr
	}

	return accessToken, nil
}
