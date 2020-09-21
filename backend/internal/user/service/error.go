package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func err2code(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func jsonDecodeValidate(db io.ReadCloser, v interface{}) error {
	err := json.NewDecoder(db).Decode(v)
	if err != nil {
		return err
	}
	type Validator interface {
		Validate() error
	}
	validable, ok := v.(Validator)
	if ok {
		err = validable.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
