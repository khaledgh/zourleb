// Package validator wires go-playground/validator into Echo with custom rules.
package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/pkg/phone"
	"github.com/zourleb/zourleb-api/pkg/response"
)

// Validator adapts go-playground/validator to Echo's Validator interface.
type Validator struct {
	v *validator.Validate
}

// New constructs the validator and registers custom rules (e.g. lb_phone).
func New() *Validator {
	v := validator.New()
	_ = v.RegisterValidation("lb_phone", func(fl validator.FieldLevel) bool {
		return phone.IsValid(fl.Field().String())
	})
	return &Validator{v: v}
}

// Validate runs struct validation. On failure it returns a localized
// VALIDATION_FAILED domain error with field-level details.
func (vv *Validator) Validate(i interface{}) error {
	if err := vv.v.Struct(i); err != nil {
		fields := map[string]string{}
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				fields[jsonName(fe.Field())] = ruleMessage(fe)
			}
		}
		return response.ErrValidation.WithFields(fields)
	}
	return nil
}

// Install registers the validator on an Echo instance.
func Install(e *echo.Echo) {
	e.Validator = New()
}

func jsonName(field string) string {
	// validator reports the Go field name; lowercasing is a pragmatic default.
	if field == "" {
		return field
	}
	return string(field[0]|0x20) + field[1:]
}

func ruleMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Must be a valid email address."
	case "lb_phone":
		return "Must be a valid Lebanese phone number."
	case "min":
		return "Value is too short."
	case "max":
		return "Value is too long."
	case "oneof":
		return "Value is not allowed."
	default:
		return "Invalid value."
	}
}
