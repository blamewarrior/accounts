package blamewarrior

import (
	"fmt"
)

type Validator struct {
	messages []string
}

func (v *Validator) MustNotBeEmpty(value string, msgArgs ...interface{}) bool {
	if value != "" {
		return true
	}

	var msg string

	if len(msgArgs) == 0 || msgArgs == nil {
		msg = fmt.Sprintf("must not be empty")
	}

	if len(msgArgs) == 1 {
		msg = msgArgs[0].(string)
	}
	if len(msgArgs) > 1 {
		msg = fmt.Sprintf(msgArgs[0].(string), msgArgs[1:]...)
	}

	v.messages = append(v.messages, msg)

	return false
}

func (v *Validator) ErrorMessages() []string {
	return v.messages
}

func (v *Validator) IsValid() bool {
	return len(v.messages) == 0
}
