package errorz

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func NewValidationError(errs validator.ValidationErrors) Response {
	var fields []Field
	for _, err := range errs {
		fields = append(fields, newValidationErrorField(err))
	}
	return Response{ErrorFields: fields}
}

func newValidationErrorField(err validator.FieldError) Field {
	return Field{
		Category: InvalidRequestError,
		Code:     fieldErrorCode(err),
		DocURL:   APIDocURL,
		Message:  fieldErrorMsg(err),
	}
}

func fieldErrorCode(err validator.FieldError) Code {
	switch err.ActualTag() {
	case "required":
		return MissingRequiredParameter
	default:
		return InvalidValue
	}
}

func fieldErrorMsg(err validator.FieldError) string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + err.Field() + "'")
	sb.WriteString(", condition: " + err.ActualTag())

	// Print condition parameters
	if err.Param() != "" {
		sb.WriteString(" { " + err.Param() + " }")
	}

	if err.Value() != nil && err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", err.Value()))
	}

	return sb.String()
}

