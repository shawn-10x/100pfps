package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/utils"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func SetupValidators(e *echo.Echo) {
	v := new(Validator)
	v.validator = validator.New(validator.WithRequiredStructEnabled())
	v.validator.RegisterValidation("tags", ValidateTags)
	v.validator.RegisterValidation("tags_max_count", TagsMaxCount)
	v.validator.RegisterValidation("tag_length", TagsLength)

	e.Validator = v
}

func ValidationErrors(v_errors validator.ValidationErrors) utils.Ms {
	errors := map[string]string{}
	errors["kind"] = "Error validating"
	for _, err := range v_errors {
		var msg string
		switch err.Tag() {
		case "required":
			msg = "Missing field"
		case "min":
			msg = "Too short"
		case "max":
			msg = "Too long"
		case "tags":
			msg = "Invalid tags"
		case "tags_max_count":
			msg = "Too many tags"
		case "tag_length":
			msg = "Tag maximum length is " + err.Param()
		default:
			msg = "Invalid"
			// fmt.Println("namespace", err.Namespace())
			// fmt.Println("field", err.Field())
			// fmt.Println("struct namespace", err.StructNamespace())
			// fmt.Println("struct field", err.StructField())
			// fmt.Println("tag", err.Tag())
			// fmt.Println("actual tag", err.ActualTag())
			// fmt.Println("kind", err.Kind())
			// fmt.Println("type", err.Type())
			// fmt.Println("value", err.Value())
			// fmt.Println("param", err.Param())
			// fmt.Println()
		}
		errors[err.Field()] = msg
	}
	return errors
}
