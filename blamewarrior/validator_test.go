package blamewarrior_test

import (
	"testing"

	"github.com/blamewarrior/users/blamewarrior"
	"github.com/stretchr/testify/assert"
)

func TestMustNotBeEmpty(t *testing.T) {
	examples := map[string]struct {
		Opts    []interface{}
		Message string
	}{
		"empty": {
			Opts:    []interface{}{},
			Message: "must not be empty",
		},
		"only 1 message": {
			Opts:    []interface{}{"title must not be empty"},
			Message: "title must not be empty",
		},
		"parameterized args": {
			Opts:    []interface{}{"title %s", "must not be empty"},
			Message: "title must not be empty",
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			v := new(blamewarrior.Validator)

			args := example.Opts

			v.MustNotBeEmpty("", args...)

			result := v.ErrorMessages()[0]

			assert.Equal(t, example.Message, result)
		})
	}
}
