package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type DeviceStatePayload struct {
	DeviceID string   `json:"device_id" validate:"required"`
	Type     string   `json:"type" validate:"required,oneof=capability property"` // "capability" or "property"
	Instance string   `json:"instance" validate:"required"`
	Value    *float32 `json:"value" validate:"required,gte=0,lte=100"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Field string
	Tag   string
	Param string
}

func validatorErrorString(val_errors validator.ValidationErrors) string {
	var err_msg string
	for idx, e := range val_errors {
		err_msg += ValidationError{
			Field: e.Field(),
			Tag:   e.Tag(),
			Param: e.Param(),
		}.toString()
		if idx != len(val_errors)-1 {
			err_msg += ", "
		}
	}
	return err_msg
}

func (e ValidationError) toString() string {
	err_msg := fmt.Sprintf("field '%s' ", e.Field)
	log.Printf("Kind: %s %s %s", e.Tag, e.Field, e.Param)
	switch e.Tag {
	case "required":
		err_msg += "is required"
	case "lte", "gte":
		err_msg += fmt.Sprintf("needs to be %s %s", e.Tag, e.Param)
	case "oneof":
		err_msg += fmt.Sprintf("needs to be one of [%s]", e.Param)
	default:
		err_msg = "validation error"
	}
	return err_msg
}
