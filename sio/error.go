package sio

import (
	"fmt"
	"net/http"
)

type ServiceError interface {
	HTTPStatus() int
	Message() string
}

type serviceError struct {
	httpStatus int
	message    string
}

func (s serviceError) HTTPStatus() int {
	return s.httpStatus
}

func (s serviceError) Message() string {
	return s.message
}

func ErrUnknownPath(path string) serviceError {
	return serviceError{http.StatusNotFound, fmt.Sprintf("unknown path: %s", path)}
}

func ErrIllegalInput(err error) serviceError {
	return serviceError{http.StatusBadRequest, fmt.Sprintf("illegal input: %s", err)}
}

func ErrTimeout() serviceError {
	return serviceError{http.StatusBadRequest, "computation timeout reached"}
}

func ErrUnsupportedContentType(ct string) serviceError {
	return serviceError{http.StatusUnsupportedMediaType, fmt.Sprintf("unsupported content-type %s", ct)}
}

func ErrUnsupportedAccept(acc string) serviceError {
	return serviceError{http.StatusUnsupportedMediaType, fmt.Sprintf("unsupported accept: %s", acc)}
}

func ErrServer(err error) serviceError {
	return serviceError{http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err)}
}

func WriteError(w http.ResponseWriter, r *http.Request, err ServiceError) {
	w.WriteHeader(err.HTTPStatus())
	state := r.Context().Value(State{}).(*ComputationState)
	state.Success = false
	state.Error = err.Message()
	result := ComputationResult{State: *state}
	WriteResult(w, r, result)
}
