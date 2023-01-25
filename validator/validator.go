package validator

import (
	"github.com/forPelevin/gomoji"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator(i *do.Injector) (echo.Validator, error) {
	v := &CustomValidator{
		Validator: validator.New(),
	}

	// Register custom validation
	v.MustRegisterValidation("exclude_emoji", ExcludeEmoji)

	return v, nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (cv *CustomValidator) MustRegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	err := cv.Validator.RegisterValidation(tag, fn, callValidationEvenIfNull...)
	if err != nil {
		panic(err)
	}
}

func ExcludeEmoji(fl validator.FieldLevel) bool {
	field := fl.Field()
	str := field.String()
	return !gomoji.ContainsEmoji(str)
}
