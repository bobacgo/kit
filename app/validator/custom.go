package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

// 自定义验证器

func registerValidation(valid *validator.Validate) {
	valid.RegisterValidation("duration", durationValid)
}

func durationValid(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "" {
		return true
	}
	_, err := time.ParseDuration(str)
	return err == nil
}