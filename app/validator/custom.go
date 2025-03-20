package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// 自定义验证器
func registerValidation(valid *validator.Validate) {
	// “validate:duration”
	// 验证时间格式 "300ms", "-1.5h" or "2h45m"
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
