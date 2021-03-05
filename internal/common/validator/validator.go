package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type defaultValidator struct {
	validate *validator.Validate
}

func NewDefaultValidator() binding.StructValidator {
	return &defaultValidator{validate: newValidate()}
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		err := v.validate.Struct(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	return v.validate
}

func newValidate() *validator.Validate {
	validate := validator.New()
	// using Binding tags as validation identifier
	validate.SetTagName("binding")
	// use JSON tags as field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})
	return validate
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}
