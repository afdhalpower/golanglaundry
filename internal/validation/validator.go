package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: '%s' validation failed on '%s'", v.Field, v.Value, v.Tag)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = fmt.Sprintf("Field '%s' %s", err.Field, tagMessage(err.Tag))
	}
	return result
}

func tagMessage(tag string) string {
	messages := map[string]string{
		"required": "wajib diisi",
		"email":    "format email tidak valid",
		"min":      "terlalu pendek",
		"max":      "terlalu panjang",
		"gt":       "harus lebih besar dari 0",
		"oneof":    "nilai tidak valid",
	}
	if msg, ok := messages[tag]; ok {
		return msg
	}
	return fmt.Sprintf("tidak valid (aturan: %s)", tag)
}

func ValidateStruct(s interface{}) ValidationErrors {
	var errs ValidationErrors

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: fmt.Sprintf("%v", err.Value()),
		})
	}

	return errs
}
