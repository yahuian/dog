package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	enLocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// https://github.com/go-playground/validator/tree/v9/_examples/gin-upgrading-overriding

type defaultValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

var instance *defaultValidator

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := v.validate.Struct(obj); err != nil {
			return errors.New(v.pretty(err))
		}
	}

	return nil
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}

func (v *defaultValidator) Engine() interface{} {
	return v.validate
}

func Init() error {
	v := &defaultValidator{
		validate: validator.New(),
	}

	// register trans
	english := enLocales.New()
	uni := ut.New(english, english)
	v.translator, _ = uni.GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(v.validate, v.translator); err != nil {
		return fmt.Errorf("validate register translation err: %w", err)
	}

	v.validate.SetTagName("validate")

	// replace gin default validate
	binding.Validator = v

	instance = v

	return nil
}

func Struct(s any) error {
	err := instance.validate.Struct(s)
	if err != nil {
		return errors.New(instance.pretty(err))
	}

	return nil
}

func Var(s any, name, tag string) error {
	err := instance.validate.Var(s, tag)
	if err != nil {
		return errors.New(name + instance.pretty(err))
	}

	return nil
}

func (v *defaultValidator) pretty(err error) string {
	var errorsMap validator.ValidationErrors

	if !errors.As(err, &errorsMap) {
		return err.Error()
	}

	result := make([]string, 0, len(errorsMap))

	for _, e := range errorsMap {
		result = append(result, e.Translate(v.translator))
	}

	return strings.Join(result, ", ")
}
