package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/major75/online-subscriptions/pkg/types"
)

var validate *validator.Validate

// GetInstance returns singleton-validator with registered custom validation rules
func GetInstance() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		_ = validate.RegisterValidation("datetime", validateCustomDateTime)
	}
	return validate
}

// validateCustomDateTime - custom validator for time.Time with format mm-yyyy
func validateCustomDateTime(fl validator.FieldLevel) bool {
	param := fl.Param()
	value := fl.Field().Interface()

	// If the field is empty consider it invalid (use omitempty for optional fields)
	if value == nil || (value == time.Time{}) {
		return false
	}

	t, ok := value.(time.Time)
	if !ok {
		return false
	}

	// Default format
	if param == "" {
		param = types.DATE_FORMAT
	}

	// Check that the date matches the format
	formatted := t.Format(param)
	_, err := time.Parse(param, formatted)
	return err == nil
}

func Validate(i interface{}) error {
	return GetInstance().Struct(i)
}
