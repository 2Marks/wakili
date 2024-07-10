package errs

import "net/http"

type ErrProxyRequest struct {
	err        error
	statusCode int
}

func (e *ErrProxyRequest) Error() string {
	return e.err.Error()
}

func ErrInternal(err error) *ErrProxyRequest {
	return &ErrProxyRequest{err: err, statusCode: http.StatusInternalServerError}
}
