package validate

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var trans ut.Translator

func init() {
	// Register translator to validator
	v, _ := binding.Validator.Engine().(*validator.Validate)

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
}

func ValidateError(err error) string {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return err.Error()
	}
	var list []string
	for _, e := range errs {
		list = append(list, e.Translate(trans))
	}
	return strings.Join(list, ";")
}
