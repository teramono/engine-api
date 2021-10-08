package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


func (server *APIServer) validateBody(ctx *gin.Context, obj interface{}) (validator.ValidationErrors, error) {
	// Get body bytes.
	data, err := ctx.GetRawData()
	if err != nil {
		return validator.ValidationErrors{}, err
	}

	// Deserialize struct to specified object.
	if err = json.Unmarshal(data, obj); err != nil {
		return validator.ValidationErrors{}, err
	}

	// Validate object with its struct tags.
	if err = server.Validator.Struct(obj); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return validator.ValidationErrors{}, err
		} else if errs, ok := err.(validator.ValidationErrors); ok {
			return errs, nil
		}
	}

	return validator.ValidationErrors{}, nil
}

