package rules

import (
	"fmt"
	"github.com/goravel/framework/contracts/validation"
)

type Confirmed struct {
	field string
}

// Signature The name of the rule.
func (receiver *Confirmed) Signature() string {
	return "confirmed"
}

// Passes Determine if the validation rule passes.
func (receiver *Confirmed) Passes(data validation.Data, val any, options ...any) bool {
	// Get the confirmation field name (field + "_confirmation")
	confirmationField := options[0].(string) + "_confirmation"

	// Get the confirmation value from the data
	confirmationValue, exists := data.Get(confirmationField)

	if !exists {
		return false
	}

	// Compare the values
	return fmt.Sprintf("%v", val) == fmt.Sprintf("%v", confirmationValue)
}

// Message Get the validation error message.
func (receiver *Confirmed) Message() string {
	return fmt.Sprintf("The %s confirmation does not match.", receiver.field)
}
