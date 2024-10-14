package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 500 - internal server error
func Dispatch500Error(w http.ResponseWriter, err error) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(WriteError(fmt.Sprintf("%v", err), nil))
}

// 400 - bad request
func Dispatch400Error(w http.ResponseWriter, msg string, err any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(WriteError(msg, err))
}

func Dispatch409Error(w http.ResponseWriter, msg string, err any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusConflict)
	w.Write(WriteError(msg, err))
}

// 403 - forbidden request, incase of non-authorised request
func Dispatch403Error(w http.ResponseWriter, msg string, err any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusForbidden)
	w.Write(WriteError(msg, err))
}

// 404 - not found
func Dispatch404Error(w http.ResponseWriter, msg string, err any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusNotFound)
	w.Write(WriteError(msg, err))
}

// 422 - unprocessable entity
func Dispatch422Error(w http.ResponseWriter, msg string, err any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(WriteError(msg, err))
}

// 200 - OK
func Dispatch200(w http.ResponseWriter, msg string, data any) {
	AddDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write(WriteInfo(msg, data))
}

func WriteInfo(message string, data any) []byte {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	r, err := json.Marshal(response)
	if err == nil {
		return r
	} else {
		log.Printf("err: %s", err)
	}
	return nil
}

func WriteError(message string, err interface{}) []byte {
	response := APIResponse{
		Success: false,
		Message: message,
		Data:    err,
	}
	data, err := json.Marshal(response)
	if err == nil {
		return data
	} else {
		log.Printf("err: %s", err)
	}
	return nil
}

func AddDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}
