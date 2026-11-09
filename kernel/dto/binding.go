package dto

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func BindStrictJSON(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(v)
}
