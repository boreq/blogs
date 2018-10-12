// Package api implements a framework for creating a JSON API.
package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var InternalServerError = NewError(http.StatusInternalServerError, "Internal server error.")
var BadRequestError = NewError(http.StatusBadRequest, "Bad request.")
var NotFoundError = NewError(http.StatusNotFound, "Not found.")
var UnauthorizedError = NewError(http.StatusUnauthorized, "Unauthorized.")

type Error interface {
	GetCode() int
	Error() string
}

type Response interface {
	GetCode() int
	GetData() interface{}
}

func NewError(code int, message string) Error {
	return apiError{Code: code, Message: message}
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err apiError) GetCode() int {
	return err.Code
}

func (err apiError) Error() string {
	return err.Message
}

func NewResponseOk(data interface{}) Response {
	return apiResponse{Code: http.StatusOK, Data: data}
}

func NewResponse(code int, data interface{}) Response {
	return apiResponse{Code: code, Data: data}
}

type apiResponse struct {
	Code int
	Data interface{}
}

func (r apiResponse) GetCode() int {
	return r.Code
}

func (r apiResponse) GetData() interface{} {
	return r.Data
}

type Handle func(r *http.Request, p httprouter.Params) (Response, Error)

func Call(w http.ResponseWriter, r *http.Request, p httprouter.Params, handle Handle) error {
	response, apiErr := handle(r, p)
	if apiErr != nil {
		response = NewResponse(apiErr.GetCode(), apiError{apiErr.GetCode(), apiErr.Error()})
	} else {
		if response == nil {
			response = NewResponseOk(nil)
		}
	}
	j, err := json.Marshal(response.GetData())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetCode())
	_, err = bytes.NewBuffer(j).WriteTo(w)
	return err
}

func Wrap(handle Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		Call(w, r, p, handle)
	}
}
