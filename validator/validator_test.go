package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator(t *testing.T) {
	dummyInjector := do.New()
	v, err := NewCustomValidator(dummyInjector)
	assert.NoError(t, err)

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
