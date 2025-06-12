package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Name          string
	Email         string
	Password      string
	RememberToken string `gorm:"column:remember_token"`
	orm.SoftDeletes
}
