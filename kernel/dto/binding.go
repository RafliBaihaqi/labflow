package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindStrictJSON(c *gin.Context, v interface{}) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(v)
}
