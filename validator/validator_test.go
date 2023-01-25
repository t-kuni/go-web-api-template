package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator(t *testing.T) {
	v := NewCustomValidator()

	t.Run("ExcludeEmoji", func(t *testing.T) {
		type Input struct {
			Str string `validate:"exclude_emoji"`
		}
		input := Input{Str: "Hello!😀"}

		err := v.Validate(input)
		vErr := err.(validator.ValidationErrors)

		assert.Len(t, vErr, 1)
		assert.Equal(t, vErr[0].Tag(), "exclude_emoji")

		input = Input{Str: "Hello!"}

		err = v.Validate(input)
		assert.NoError(t, err)
	})
}
