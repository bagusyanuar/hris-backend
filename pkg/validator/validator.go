package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// Validate is a global singleton for the validator.
var Validate *validator.Validate
var trans ut.Translator

func init() {
	Validate = validator.New()

	// Setup translator for Indonesian locale
	idLocale := id.New()
	uni := ut.New(idLocale, idLocale)
	trans, _ = uni.GetTranslator("id")

	// Register default Indonesian translations
	_ = id_translations.RegisterDefaultTranslations(Validate, trans)

	// Register tag name function to extract the 'json' tag name instead of the struct field name.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
}

// ValidateStruct validates a given struct and returns a map of human-readable field errors (Laravel format).
// It returns nil if there are no errors.
func ValidateStruct(data any) map[string][]string {
	err := Validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string][]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			errors[field] = append(errors[field], e.Translate(trans))
		}
	} else {
		// If it's not a ValidationErrors type, return a generic error.
		errors["general"] = []string{err.Error()}
	}
	return errors
}
