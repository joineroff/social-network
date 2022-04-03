package validation

import (
	"github.com/go-playground/validator/v10"
)

// Singleton used for caching struct tags
var (
	V = validator.New()
)

type ValidationErrors = validator.ValidationErrors

func Struct(s interface{}) error {
	return validator.New().Struct(s)
}
