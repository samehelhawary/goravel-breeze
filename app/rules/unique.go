package rules

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type Unique struct {
}

// Signature returns the name of the rule.
func (receiver *Unique) Signature() string {
	return "unique"
}

// Passes determines if the validation rule passes.
func (receiver *Unique) Passes(data validation.Data, val any, options ...any) bool {
	var isExists bool
	var tableName = options[0].(string)
	var columnName = options[1].(string)
	err := facades.Orm().Query().Table(tableName).Where(columnName, val).Exists(&isExists)
	if err != nil {
		return true
	}

	if isExists {
		return false
	}

	return true
}

// Message gets the validation error message.
func (receiver *Unique) Message() string {
	return "The :attribute has already been taken."
}
