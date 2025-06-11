package facades

import (
	"log"

	breeze "github.com/samehelhawary/goravel-breeze"
	"github.com/samehelhawary/goravel-breeze/contracts"
)

func Breeze() contracts.Breeze {
	instance, err := breeze.App.Make(breeze.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Breeze)
}
